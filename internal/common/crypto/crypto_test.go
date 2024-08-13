package crypto_test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"testing"

	"github.com/Alekseyt9/ypmetrics/internal/common/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadPrivateKey(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	privateKeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		},
	)

	fileName := "test_private_key.pem"
	err = os.WriteFile(fileName, privateKeyPEM, 0600)
	require.NoError(t, err)
	defer os.Remove(fileName)

	loadedPrivateKey, err := crypto.LoadPrivateKey(fileName)
	require.NoError(t, err)

	assert.IsType(t, &rsa.PrivateKey{}, loadedPrivateKey)
}

func TestLoadPublicKey(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	publicKeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: x509.MarshalPKCS1PublicKey(&privateKey.PublicKey),
		},
	)

	fileName := "test_public_key.pem"
	err = os.WriteFile(fileName, publicKeyPEM, 0600)
	require.NoError(t, err)
	defer os.Remove(fileName)

	loadedPublicKey, err := crypto.LoadPublicKey(fileName)
	require.NoError(t, err)

	assert.IsType(t, &rsa.PublicKey{}, loadedPublicKey)
}

func TestCipherAndDecipher(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	publicKey := &privateKey.PublicKey

	originalData := []byte("Test data for encryptionTest data for encryptionTest data for encryptionTest data for encryptionTest data for encryptionTest data for encryptionTest data for encryptionTest data for encryptionTest data for encryptionTest data for encryptionTest data for encryptionTest data for encryptionTest data for encryptionTest data for encryptionTest data for encryptionTest data for encryptionTest data for encryptionTest data for encryptionTest data for encryptionTest data for encryptionTest data for encryptionTest data for encryptionTest data for encryptionTest data for encryptionTest data for encryptionTest data for encryptionTest data for encryptionTest data for encryptionTest data for encryptionTest data for encryptionTest data for encryptionTest data for encryptionTest data for encryptionTest data for encryptionTest data for encryptionTest data for encryptionTest data for encryptionTest data for encryptionTest data for encryptionTest data for encryptionTest data for encryptionTest data for encryption")

	encryptedData, err := crypto.Cipher(originalData, publicKey)
	require.NoError(t, err)

	assert.NotNil(t, encryptedData)
	assert.NotEqual(t, originalData, encryptedData)

	decryptedData, err := crypto.Decipher(encryptedData, privateKey)
	require.NoError(t, err)

	assert.NotNil(t, decryptedData)
	assert.Equal(t, originalData, decryptedData)
}
