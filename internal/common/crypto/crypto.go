package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"os"
)

func LoadPrivateKey(fileName string) (*rsa.PrivateKey, error) {
	pKeyBytes, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	pKeyBlock, _ := pem.Decode(pKeyBytes)

	pKey, err := x509.ParsePKCS1PrivateKey(pKeyBlock.Bytes)
	if err != nil {
		return nil, err
	}

	return pKey, nil
}

func LoadPublicKey(fileName string) (*rsa.PublicKey, error) {
	pKeyBytes, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	pKeyBlock, _ := pem.Decode(pKeyBytes)

	pKey, err := x509.ParsePKCS1PublicKey(pKeyBlock.Bytes)
	if err != nil {
		return nil, err
	}

	return pKey, nil
}

func Cipher(data []byte, rsaPubKey *rsa.PublicKey) ([]byte, error) {
	maxBlockSize := rsaPubKey.Size() - 2*sha256.Size - 2

	var encryptedData []byte
	for start := 0; start < len(data); start += maxBlockSize {
		end := start + maxBlockSize
		if end > len(data) {
			end = len(data)
		}

		block := data[start:end]
		encryptedBlock, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, rsaPubKey, block, nil)
		if err != nil {
			return nil, err
		}

		encryptedData = append(encryptedData, encryptedBlock...)
	}

	return encryptedData, nil
}

func Decipher(data []byte, key *rsa.PrivateKey) ([]byte, error) {
	rsaKeySize := key.Size()

	var decryptedData []byte
	for start := 0; start < len(data); start += rsaKeySize {
		end := start + rsaKeySize
		if end > len(data) {
			end = len(data)
		}

		block := data[start:end]
		decryptedBlock, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, key, block, nil)
		if err != nil {
			return nil, err
		}

		decryptedData = append(decryptedData, decryptedBlock...)
	}

	return decryptedData, nil
}
