package crypto_test

import (
	"testing"
	"vernam_cipher/crypto"

	"github.com/stretchr/testify/assert"
)

func TestGenerateKey(t *testing.T) {
	t.Parallel()

	message := "Hello world!"
	messageLength := len(message)

	key, err := crypto.GenerateKey(messageLength)

	assert.Nil(t, err)
	assert.Equal(t, messageLength, len(key))
}

func TestEncrypt(t *testing.T) {
	t.Parallel()

	message := []byte("Hello world!")
	key, _ := crypto.GenerateKey(len(message))

	encrypted, err := crypto.Encrypt(message, key)

	assert.Nil(t, err)
	assert.Equal(t, len(message), len(encrypted))
	assert.NotEqual(t, message, encrypted)
}

func TestDecrypt(t *testing.T) {
	t.Parallel()

	message := []byte("Hello world!")
	key, _ := crypto.GenerateKey(len(message))
	encrypted, _ := crypto.Encrypt(message, key)

	decrypted, err := crypto.Decrypt(encrypted, key)

	assert.Nil(t, err)
	assert.Equal(t, len(message), len(decrypted))
	assert.Equal(t, message, decrypted)
}
