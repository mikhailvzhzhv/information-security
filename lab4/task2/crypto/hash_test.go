package crypto_test

import (
	"fmt"
	"gost_28147-89/crypto"
	"testing"
)

func Test(t *testing.T) {
	msg := []byte("askjsallk lskdckllsad s")
	hash := crypto.MerkleDamgardHash(msg)
	fmt.Printf("Hash: %x\n", hash)

	crypto.TestAvalancheEffect()
}
