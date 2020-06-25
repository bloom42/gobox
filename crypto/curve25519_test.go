package crypto



func TestPrivateKeyEncryptDecrypt(t *testing.T) {
	message := []byte("this is a simple message")
	nonce, err := RandBytes(AEADNonceSize)
	if err != nil {
		t.Error(err)
	}

	toPublicKey, toPrivateKey, err := GenerateKeyPair(RandReader())
	if err != nil {
		t.Error(err)
	}

	fromPublicKey, fromPrivateKey, err := GenerateKeyPair(RandReader())
	if err != nil {
		t.Error(err)
	}

	ciphertext, err := toPublicKey.Encrypt(fromPrivateKey, nonce, message)
	if err != nil {
		t.Error(err)
	}

	plaintext, err := toPrivateKey.Decrypt(fromPublicKey, nonce, ciphertext)
	if err != nil {
		t.Error(err)
	}

	if !ConstantTimeCompare(message, plaintext) {
		t.Errorf("Message (%s) and plaintext (%s) don't match", string(message), string(plaintext))
	}
}

func TestPrivateKeyEncryptDecryptEphemeral(t *testing.T) {
	message := []byte("this is a simple message")

	toPublicKey, toPrivateKey, err := GenerateKeyPair(RandReader())
	if err != nil {
		t.Error(err)
	}

	ciphertext, ephemeralPublicKey, err := toPublicKey.EncryptEphemeral(message)
	if err != nil {
		t.Error(err)
	}

	plaintext, err := toPrivateKey.DecryptEphemeral(ephemeralPublicKey, ciphertext)
	if err != nil {
		t.Error(err)
	}

	if !ConstantTimeCompare(message, plaintext) {
		t.Errorf("Message (%s) and plaintext (%s) don't match", string(message), string(plaintext))
	}
}
