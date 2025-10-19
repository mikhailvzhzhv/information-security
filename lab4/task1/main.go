package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"vernam_cipher/crypto"
)

func main() {
	fmt.Println("type anything: ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	message := scanner.Text()

	fmt.Printf("message : %s\n", message)

	key, err := crypto.GenerateKey(len(message))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("generated key: %s\n", key)

	encrypted, err := crypto.Encrypt([]byte(message), key)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("encrypted message: %s\n", encrypted)

	decrypted, err := crypto.Decrypt(encrypted, key)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("decrypted message: %s\n", decrypted)
}
