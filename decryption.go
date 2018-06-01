package main

import (
	"fmt"

	openssl "github.com/Luzifer/go-openssl"
)

type decryptMethod func(body []byte, passphrase string) ([]byte, error)

func decryptMethodFromName(name string) (decryptMethod, error) {
	switch name {

	case "openssl-md5":
		return decryptOpenSSLMD5, nil

	default:
		return nil, fmt.Errorf("Decrypt method %q not found", name)
	}
}

func decryptOpenSSLMD5(body []byte, passphrase string) ([]byte, error) {
	return openssl.New().DecryptString(cfg.Password, string(body))
}
