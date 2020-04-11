// Package email provides an easy to use and hard to misuse email API
package email

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net/mail"
	"net/smtp"
	"net/textproto"
	"strings"
	"time"

	"gitlab.com/bloom42/lily/crypto"
)

const (
	// MaxLineLength is the maximum line length per RFC 2045
	MaxLineLength = 76
	// defaultContentType is the default Content-Type according to RFC 2045, section 5.2
	defaultContentType = "text/plain; charset=us-ascii"
)

// ErrMissingBoundary is returned when there is no boundary given for a multipart entity
var ErrMissingBoundary = errors.New("No boundary found for multipart entity")

// ErrMissingContentType is returned when there is no "Content-Type" header for a MIME entity
var ErrMissingContentType = errors.New("No Content-Type found for MIME entity")

// Email is an email...
// either Text or HTML must be provided
type Email struct {
	ReplyTo     []string
	From        string
	To          []string
	Bcc         []string
	Cc          []string
	Subject     string
	Text        []byte // Plaintext message
	HTML        []byte // Html message
	Headers     textproto.MIMEHeader
	Attachments []Attachment
	// ReadReceipt []string
}

func (email *Email) Bytes() ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})
	hasAttachements := len(email.Attachments) > 0
	isAlternative := len(email.Text) > 0 && len(email.HTML) > 0
	var writer *multipart.Writer

	headers, err := email.headers()
	if err != nil {
		return nil, err
	}

	if hasAttachements || isAlternative {
		writer = multipart.NewWriter(buffer)
	}
	switch {
	case hasAttachements:
		headers.Set("Content-Type", "multipart/mixed;\r\n boundary="+writer.Boundary())
	case isAlternative:
		headers.Set("Content-Type", "multipart/alternative;\r\n boundary="+writer.Boundary())
	case len(email.HTML) > 0:
		headers.Set("Content-Type", "text/html; charset=UTF-8")
		headers.Set("Content-Transfer-Encoding", "quoted-printable")
	default:
		headers.Set("Content-Type", "text/plain; charset=UTF-8")
		headers.Set("Content-Transfer-Encoding", "quoted-printable")
	}

	_, err = io.WriteString(buffer, "\r\n")
	if err != nil {
		return nil, err
	}

	// Check to see if there is a Text or HTML field
	if len(email.Text) > 0 || len(email.HTML) > 0 {
		var subWriter *multipart.Writer

		if hasAttachements && isAlternative {
			// Create the multipart alternative part
			subWriter = multipart.NewWriter(buffer)
			header := textproto.MIMEHeader{
				"Content-Type": {"multipart/alternative;\r\n boundary=" + subWriter.Boundary()},
			}
			if _, err := writer.CreatePart(header); err != nil {
				return nil, err
			}
		} else {
			subWriter = writer
		}
		// Create the body sections
		if len(email.Text) > 0 {
			// Write the text
			if err := writeMessage(buffer, email.Text, hasAttachements || isAlternative, "text/plain", subWriter); err != nil {
				return nil, err
			}
		}
		if len(email.HTML) > 0 {
			// Write the HTML
			if err := writeMessage(buffer, email.HTML, hasAttachements || isAlternative, "text/html", subWriter); err != nil {
				return nil, err
			}
		}
		if hasAttachements && isAlternative {
			if err := subWriter.Close(); err != nil {
				return nil, err
			}
		}
	}
	// Create attachment part, if necessary
	for _, a := range email.Attachments {
		ap, err := writer.CreatePart(a.Header)
		if err != nil {
			return nil, err
		}
		// Write the base64Wrapped content to the part
		base64Wrap(ap, a.Content)
	}
	if hasAttachements || isAlternative {
		if err := writer.Close(); err != nil {
			return nil, err
		}
	}
	return buffer.Bytes(), nil
}

func (email *Email) headers() (textproto.MIMEHeader, error) {
	res := make(textproto.MIMEHeader, len(email.Headers)+6)
	if email.Headers != nil {
		for _, h := range []string{"Reply-To", "To", "Cc", "From", "Subject", "Date", "Message-Id", "MIME-Version"} {
			if v, ok := email.Headers[h]; ok {
				res[h] = v
			}
		}
	}
	// Set headers if there are values.
	if _, ok := res["Reply-To"]; !ok && len(email.ReplyTo) > 0 {
		res.Set("Reply-To", strings.Join(email.ReplyTo, ", "))
	}
	if _, ok := res["To"]; !ok && len(email.To) > 0 {
		res.Set("To", strings.Join(email.To, ", "))
	}
	if _, ok := res["Cc"]; !ok && len(email.Cc) > 0 {
		res.Set("Cc", strings.Join(email.Cc, ", "))
	}
	if _, ok := res["Subject"]; !ok && email.Subject != "" {
		res.Set("Subject", email.Subject)
	}
	if _, ok := res["Message-Id"]; !ok {
		id, err := generateMessageID()
		if err != nil {
			return nil, err
		}
		res.Set("Message-Id", id)
	}
	// Date and From are required headers.
	if _, ok := res["From"]; !ok {
		res.Set("From", email.From)
	}
	if _, ok := res["Date"]; !ok {
		res.Set("Date", time.Now().Format(time.RFC1123Z))
	}
	if _, ok := res["MIME-Version"]; !ok {
		res.Set("MIME-Version", "1.0")
	}
	for field, vals := range email.Headers {
		if _, ok := res[field]; !ok {
			res[field] = vals
		}
	}
	return res, nil
}

// Attachment is an email attachment.
// Based on the mime/multipart.FileHeader struct, Attachment contains the name, MIMEHeader, and content of the attachment in question
type Attachment struct {
	Filename string
	Header   textproto.MIMEHeader
	Content  []byte
}

// Send an Email
func Send(email Email, smtpHost string, smtpPort uint16, smtpAuth smtp.Auth) error {
	smtpAddress := fmt.Sprintf("%s:%d", smtpHost, smtpPort)

	// Merge the To, Cc, and Bcc fields
	to := make([]string, len(email.To)+len(email.Cc)+len(email.Bcc))
	to = append(to, email.To...)
	to = append(to, email.Bcc...)
	to = append(to, email.Cc...)
	for i, recipient := range to {
		recipientAddress, err := mail.ParseAddress(recipient)
		if err != nil {
			return fmt.Errorf("email: Invalid recipient address %v", err)
		}
		to[i] = recipientAddress.Address
	}

	// Check to make sure there is at least one recipient and one "From" address
	if email.From == "" || len(to) == 0 {
		return errors.New("email: Must specify at least one From address and one To address")
	}

	from, err := mail.ParseAddress(email.From)
	if err != nil {
		return fmt.Errorf("email: Invalid From address %v", err)
	}

	rawEmail, err := email.Bytes()
	if err != nil {
		return err
	}

	return smtp.SendMail(smtpAddress, smtpAuth, from.Address, to, rawEmail)
}

// headerToBytes renders "header" to "buff". If there are multiple values for a
// field, multiple "Field: value\r\n" lines will be emitted.
func headerToBytes(buff io.Writer, header textproto.MIMEHeader) {
	for field, vals := range header {
		for _, subval := range vals {
			// bytes.Buffer.Write() never returns an error.
			io.WriteString(buff, field)
			io.WriteString(buff, ": ")
			// Write the encoded header if needed
			switch {
			case field == "Content-Type" || field == "Content-Disposition":
				buff.Write([]byte(subval))
			default:
				buff.Write([]byte(mime.QEncoding.Encode("UTF-8", subval)))
			}
			io.WriteString(buff, "\r\n")
		}
	}
}

// generateMessageID generates and returns a string suitable for an RFC 2822
// compliant Message-ID, email.g.:
// <1444789264909237300.3464.1819418242800517193@DESKTOP01>
//
// The following parameters are used to generate a Message-ID:
// - The nanoseconds since Epoch
// - The calling PID
// - A cryptographically random int64
// - The sending hostname
func generateMessageID() (string, error) {
	t := time.Now().UnixNano()
	pid, err := crypto.RandInt64(0, 999)
	if err != nil {
		return "", err
	}
	rint, err := crypto.RandInt64(0, 999)
	if err != nil {
		return "", err
	}
	if err != nil {
		return "", err
	}
	hostname := "localhost.localdomain"
	msgid := fmt.Sprintf("<%d.%d.%d@%s>", t, pid, rint, hostname)
	return msgid, nil
}

func writeMessage(buffer io.Writer, msg []byte, multipart bool, mediaType string, w *multipart.Writer) error {
	if multipart {
		header := textproto.MIMEHeader{
			"Content-Type":              {mediaType + "; charset=UTF-8"},
			"Content-Transfer-Encoding": {"quoted-printable"},
		}
		if _, err := w.CreatePart(header); err != nil {
			return err
		}
	}

	qp := quotedprintable.NewWriter(buffer)
	// Write the text
	if _, err := qp.Write(msg); err != nil {
		return err
	}
	return qp.Close()
}

// base64Wrap encodes the attachment content, and wraps it according to RFC 2045 standards (every 76 chars)
// The output is then written to the specified io.Writer
func base64Wrap(writer io.Writer, b []byte) {
	// 57 raw bytes per 76-byte base64 linemail.
	const maxRaw = 57
	// Buffer for each line, including trailing CRLF.
	buffer := make([]byte, MaxLineLength+len("\r\n"))
	copy(buffer[MaxLineLength:], "\r\n")
	// Process raw chunks until there's no longer enough to fill a linemail.
	for len(b) >= maxRaw {
		base64.StdEncoding.Encode(buffer, b[:maxRaw])
		writer.Write(buffer)
		b = b[maxRaw:]
	}
	// Handle the last chunk of bytes.
	if len(b) > 0 {
		out := buffer[:base64.StdEncoding.EncodedLen(len(b))]
		base64.StdEncoding.Encode(out, b)
		out = append(out, "\r\n"...)
		writer.Write(out)
	}
}
