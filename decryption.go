package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	openssl "github.com/Luzifer/go-openssl"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
)

type decryptMethod func(body []byte, passphrase string) ([]byte, error)

func decryptMethodFromName(name string) (decryptMethod, error) {
	switch name {

	case "gpg-symmetric":
		return decryptGPGSymmetric, nil

	case "openssl-md5":
		return decryptOpenSSLMD5, nil

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

	md, err := openpgp.ReadMessage(msgReader, nil, func(keys []openpgp.Key, symmetric bool) ([]byte, error) {
		return []byte(passphrase), nil
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("Unable to decrypt message: %s", err)
	}

	return ioutil.ReadAll(md.UnverifiedBody)
}

func decryptOpenSSLMD5(body []byte, passphrase string) ([]byte, error) {
	return openssl.New().DecryptString(cfg.Password, string(body))
}
