package crypto_test

import (
	"fmt"
	"gost_28147-89/crypto"
	"testing"
)

func Test(t *testing.T) {
	key, _ := crypto.GenerateKey()
	msg := []byte("askjsallk lskdckllsad s")
	hash := crypto.Imitovstavka(msg, key)
	fmt.Printf("Hash: %x\n", hash)

	crypto.TestAvalancheEffect()
}
