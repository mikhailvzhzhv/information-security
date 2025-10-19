package main

import (
	"fmt"
	"gost_28147-89/crypto"
)

func main() {
	message := []byte("Bimmo world!! ImStud@22.ru")

	key, _ := crypto.GenerateKey()
	block := crypto.NewBlock(message)

	fmt.Printf("message: %s\n", message)

	crypto.EncryptBlock(block, key)

	fmt.Printf("encrypted: %s%s\n", block.L, block.R)

	crypto.Decrypt(block, key)

	fmt.Printf("decrypted: %s%s\n", block.L, block.R)
}
