// Package email provides an easy to use and hard to misuse email API
package email

import (
	"net/textproto"
)

// Email is an email...
// either Text or HTML must be provided
type Email struct {
	ReplyTo     []string
	From        string
	To          []string
	Bcc         []string
	Cc          []string
	Subject     string
	Text        []byte  // Plaintext message
	HTML        []byte  // Html message
	Sender      *string // override From as SMTP envelope sender
	Headers     textproto.MIMEHeader
	Attachments []Attachment
	// ReadReceipt []string
}

// Attachment is an email attachment.
// Based on the mime/multipart.FileHeader struct, Attachment contains the name, MIMEHeader, and content of the attachment in question
type Attachment struct {
	Filename string
	Header   textproto.MIMEHeader
	Content  []byte
}

// Send an Email
func Send(email Email) error {
	return nil
}
