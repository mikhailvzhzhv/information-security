package crypto_test

import (
	"gost_28147-89/crypto"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecrypt(t *testing.T) {
	message := []byte("Test alg0rIthm!!111")
	key, _ := crypto.GenerateKey()
	blocks := crypto.Encrypt(message, key)
	decrypted := crypto.Decrypt(blocks, key)

	assert.Equal(t, message, decrypted)
}
