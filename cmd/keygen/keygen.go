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
	GenKeys("private.key", "public.key")
}

func GenKeys(privFile, pubFile string) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Fatal(err)
	}

	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	privateKeyFile, err := os.Create(privFile)
	if err != nil {
		log.Fatal(err)
	}

	err = pem.Encode(privateKeyFile, privateKeyPEM)
	if err != nil {
		log.Fatal(err)
	}

	publicKeyPEM := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(&privateKey.PublicKey),
	}

	publicKeyFile, err := os.Create(pubFile)
	if err != nil {
		log.Fatal(err)
	}

	err = pem.Encode(publicKeyFile, publicKeyPEM)
	if err != nil {
		log.Fatal(err)
	}
}
