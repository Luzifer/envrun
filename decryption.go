package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"

	openssl "github.com/Luzifer/go-openssl/v3"
)

type decryptMethod func(body []byte, passphrase string) ([]byte, error)

func decryptMethodFromName(name string) (decryptMethod, error) {
	switch name {

	case "gpg-symmetric":
		return decryptGPGSymmetric, nil

	case "openssl-md5":
		return decryptOpenSSL(openssl.DigestMD5Sum), nil

	case "openssl-sha256":
		return decryptOpenSSL(openssl.DigestSHA256Sum), nil

	default:
		return nil, fmt.Errorf("Decrypt method %q not found", name)
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
		return nil, fmt.Errorf("Unable to read armor: %s", err)
	}

	var passwordRetry bool
	md, err := openpgp.ReadMessage(msgReader, nil, func(keys []openpgp.Key, symmetric bool) ([]byte, error) {
		if passwordRetry {
			return nil, fmt.Errorf("Wrong passphrase supplied")
		}

		passwordRetry = true
		return []byte(passphrase), nil
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("Unable to decrypt message: %s", err)
	}

	return ioutil.ReadAll(md.UnverifiedBody)
}

func decryptOpenSSL(kdf openssl.DigestFunc) decryptMethod {
	return func(body []byte, passphrase string) ([]byte, error) {
		return openssl.New().DecryptBytes(cfg.Password, body, kdf)
	}
}
