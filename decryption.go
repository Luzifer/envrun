package main

import (
	"bytes"
	"fmt"
	"io"

	openssl "github.com/Luzifer/go-openssl/v4"
	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/armor"
)

type decryptMethod func(body []byte, passphrase string) ([]byte, error)

func decryptMethodFromName(name string) (decryptMethod, error) {
	switch name {
	case "gpg-symmetric":
		return decryptGPGSymmetric, nil

	case "openssl-md5":
		return decryptOpenSSL(openssl.BytesToKeyMD5), nil

	case "openssl-sha256":
		return decryptOpenSSL(openssl.BytesToKeySHA256), nil

	default:
		return nil, fmt.Errorf("decrypt method %q not found", name)
	}
}

func decryptGPGSymmetric(body []byte, passphrase string) ([]byte, error) {
	var msgReader io.Reader

	block, err := armor.Decode(bytes.NewReader(body))
	switch err {
	case nil:
		msgReader = block.Body
	case io.EOF:
		msgReader = bytes.NewReader(body)
	default:
		return nil, fmt.Errorf("reading armor: %w", err)
	}

	var passwordRetry bool
	md, err := openpgp.ReadMessage(msgReader, nil, func(_ []openpgp.Key, _ bool) ([]byte, error) {
		if passwordRetry {
			return nil, fmt.Errorf("wrong passphrase supplied")
		}

		passwordRetry = true
		return []byte(passphrase), nil
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("decrypting message: %w", err)
	}

	data, err := io.ReadAll(md.UnverifiedBody)
	if err != nil {
		return nil, fmt.Errorf("reading GPG body: %w", err)
	}
	return data, nil
}

func decryptOpenSSL(kdf openssl.CredsGenerator) decryptMethod {
	return func(body []byte, _ string) ([]byte, error) {
		data, err := openssl.New().DecryptBytes(cfg.Password, body, kdf)
		if err != nil {
			return nil, fmt.Errorf("decrypting data: %w", err)
		}
		return data, nil
	}
}
