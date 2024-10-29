package utils

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)


func ReadPublicPEM(path string) (*rsa.PublicKey, error) {
	publicKeyPEM, err := os.ReadFile(path)
    if err != nil {
		return nil, err
    }
    publicKeyBlock, _ := pem.Decode(publicKeyPEM)
    publicKey, err := x509.ParsePKCS1PublicKey(publicKeyBlock.Bytes)
    if err != nil {
		return nil, err
    }

	return publicKey, nil
}

func ReadPrivatePEM(path string) (*rsa.PrivateKey, error) {
	privateKeyPEM, err := os.ReadFile(path)
    if err != nil {
		return nil, err
    }
	privateKeyBlock, _ := pem.Decode(privateKeyPEM)
	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}
