package main

import (
	"crypto/x509"
	"encoding/pem"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenKeys(t *testing.T) {
	privFile := "test_private.key"
	pubFile := "test_public.key"

	GenKeys(privFile, pubFile)

	_, err := os.Stat(privFile)
	require.NoError(t, err, "Private key file was not created")

	_, err = os.Stat(pubFile)
	require.NoError(t, err, "Public key file was not created")

	privData, err := os.ReadFile(privFile)
	require.NoError(t, err, "Error reading private key file")

	privBlock, _ := pem.Decode(privData)
	require.NotNil(t, privBlock, "Failed to decode PEM block for private key")
	assert.Equal(t, "RSA PRIVATE KEY", privBlock.Type, "Invalid PEM format for private key")

	_, err = x509.ParsePKCS1PrivateKey(privBlock.Bytes)
	require.NoError(t, err, "Error parsing private key")

	pubData, err := os.ReadFile(pubFile)
	require.NoError(t, err, "Error reading public key file")

	pubBlock, _ := pem.Decode(pubData)
	require.NotNil(t, pubBlock, "Failed to decode PEM block for public key")
	assert.Equal(t, "PUBLIC KEY", pubBlock.Type, "Invalid PEM format for public key")

	_, err = x509.ParsePKCS1PublicKey(pubBlock.Bytes)
	require.NoError(t, err, "Error parsing public key")

	err = os.Remove(privFile)
	require.NoError(t, err, "Error removing private key file")

	err = os.Remove(pubFile)
	require.NoError(t, err, "Error removing public key file")
}
