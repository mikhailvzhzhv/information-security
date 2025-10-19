package crypto

import (
	"crypto/rand"
	"fmt"
)

func Encrypt(message []byte, key []byte) ([]byte, error) {
	messageLength := len(message)
	if messageLength != len(key) {
		return nil, fmt.Errorf("message length [%d] not equal key length [%d]", messageLength, len(key))
	}

	encrypted := make([]byte, messageLength)
	for i := range messageLength {
		encrypted[i] = message[i] ^ key[i]
	}

	return encrypted, nil
}

func Decrypt(encrypted []byte, key []byte) ([]byte, error) {
	encryptedLength := len(encrypted)
	if len(encrypted) != len(key) {
		return nil, fmt.Errorf("encrypted length [%d] not equal key length [%d]", encryptedLength, len(key))
	}

	message := make([]byte, encryptedLength)
	for i := range len(encrypted) {
		message[i] = encrypted[i] ^ key[i]
	}

	return message, nil
}

func GenerateKey(messageLength int) ([]byte, error) {
	if messageLength < 0 {
		return nil, fmt.Errorf("message length less than 0")
	}

	key := make([]byte, messageLength)

	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}

	return key, nil
}
