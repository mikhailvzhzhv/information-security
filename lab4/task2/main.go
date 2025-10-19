package main

import (
	"fmt"
	"gost_28147-89/crypto"
	"log"
	"os"
	"path/filepath"
)

func main() {
	var filepathStr string

	fmt.Print("enter filepath: ")
	fmt.Scan(&filepathStr)

	data, err := os.ReadFile(filepathStr)
	if err != nil {
		log.Fatal(err)
	}

	key, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	base := filepath.Base(filepathStr)
	ext := filepath.Ext(base)
	nameWithoutExt := base[:len(base)-len(ext)]

	encryptedFilepath := nameWithoutExt + "_encrypted"
	decryptedFilepath := nameWithoutExt + "_decrypted" + ext

	encrypted := crypto.Encrypt(data, key)
	efile, err := os.Create(encryptedFilepath)
	if err != nil {
		log.Fatal(err)
	}

	_, err = efile.Write(encrypted)
	if err != nil {
		efile.Close()
		log.Fatal(err)
	}
	fmt.Println("Encrypted file saved as:", encryptedFilepath)

	data, err = os.ReadFile(encryptedFilepath)
	if err != nil {
		log.Fatal(err)
	}

	decrypted := crypto.Decrypt(data, key)
	dfile, err := os.Create(decryptedFilepath)
	if err != nil {
		log.Fatal(err)
	}
	defer dfile.Close()

	_, err = dfile.Write(decrypted)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Decrypted file saved as:", decryptedFilepath)
}
