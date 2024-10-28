package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"
)


func main() {

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		log.Fatal(err.Error)
	}

	publicKey := &privateKey.PublicKey

	pemPrivateFile, err := os.Create("private_key.pem")
	pemPublicFile, err := os.Create("public_key.pem")
	defer pemPrivateFile.Close()

	if err != nil {
		log.Fatal(err)
	}

	var pemPrivateBlock = &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	publicKeyBytes := x509.MarshalPKCS1PublicKey(publicKey)
	if err != nil {
		log.Fatal(err)
	}

	var pemPublicBlock = &pem.Block{
        Type:  "RSA PUBLIC KEY",
        Bytes: publicKeyBytes,
    }

	err = pem.Encode(pemPrivateFile, pemPrivateBlock)

	if err != nil {
		log.Fatal(err)
	}

	err = pem.Encode(pemPublicFile, pemPublicBlock)

	if err != nil {
		log.Fatal(err)
	}


}
