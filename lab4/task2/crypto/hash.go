package crypto

import (
	"fmt"
	"math/bits"
	"math/rand"
)

func MerkleDamgardHash(message []byte) []byte {
	H := []byte("abcdefgh")
	blocks := NewBlocks(message)

	for _, block := range blocks.blocks {
		var blockb []byte
		blockb = append(blockb, block.L[:]...)
		blockb = append(blockb, block.R[:]...)

		var keyb []byte
		keyb = append(keyb, blockb...)
		keyb = append(keyb, blockb...)
		keyb = append(keyb, blockb...)
		keyb = append(keyb, blockb...)

		key := &Key{[32]byte(keyb)}

		newBlock := NewBlock(H)
		EncryptBlock(newBlock, key)

		newH := newBlock.ToBytes()

		for i := 0; i < 8; i++ {
			H[i] ^= newH[i]
		}
	}

	return H
}

func HammingDistance(a, b []byte) int {
	if len(a) != len(b) {
		panic("hash lengths differ")
	}
	diff := 0
	for i := range a {
		diff += bits.OnesCount8(a[i] ^ b[i])
	}
	return diff
}

func TestAvalancheEffect() {
	message := []byte("The quick brown fox jumps over the lazy dog")

	hash1 := MerkleDamgardHash(message)

	message2 := make([]byte, len(message))
	copy(message2, message)

	byteIndex := rand.Intn(len(message2))
	bitIndex := rand.Intn(8)
	message2[byteIndex] ^= 1 << bitIndex

	fmt.Println(string(message))
	fmt.Println(string(message2))

	hash2 := MerkleDamgardHash(message2)

	dist := HammingDistance(hash1, hash2)
	totalBits := len(hash1) * 8
	ratio := float64(dist) / float64(totalBits)

	fmt.Printf("Исходное сообщение: %q\n", message)
	fmt.Printf("Изменённый байт #%d, бит #%d\n", byteIndex, bitIndex)
	fmt.Printf("Хеш 1: %x\n", hash1)
	fmt.Printf("Хеш 2: %x\n", hash2)
	fmt.Printf("Изменилось %d бит из %d (%.2f%%)\n", dist, totalBits, ratio*100)
}
