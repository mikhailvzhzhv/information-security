package crypto

import (
	"fmt"
	"math/bits"
	"math/rand"
)

func Imitovstavka(message []byte, key *Key) []byte {
	blocks := NewBlocks(message)
	accumulator := blocks.blocks[0].ToBytes()

	for k, block := range blocks.blocks {
		blockBytes := block.ToBytes()

		encryptBlock := &Block{
			L: [4]byte{accumulator[0], accumulator[1], accumulator[2], accumulator[3]},
			R: [4]byte{accumulator[4], accumulator[5], accumulator[6], accumulator[7]},
		}

		for i := 0; i < 16; i++ {
			TransformBlock(encryptBlock, key.GetBlock(i))
		}

		if k < len(blocks.blocks)-1 {
			encryptedBytes := encryptBlock.ToBytes()
			var open []byte
			open = append(open, blocks.blocks[k+1].L[:]...)
			open = append(open, blocks.blocks[k+1].R[:]...)
			for j := 0; j < len(blockBytes); j++ {
				encryptedBytes[j] ^= open[j]
			}

			copy(accumulator, encryptedBytes)
		}
	}

	return accumulator[0:4]
}

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
	key, _ := GenerateKey()

	hash1 := Imitovstavka(message, key)

	message2 := make([]byte, len(message))
	copy(message2, message)

	byteIndex := rand.Intn(len(message2))
	bitIndex := rand.Intn(8)
	message2[byteIndex] ^= 1 << bitIndex

	fmt.Println(string(message))
	fmt.Println(string(message2))

	hash2 := Imitovstavka(message2, key)

	dist := HammingDistance(hash1, hash2)
	totalBits := len(hash1) * 8
	ratio := float64(dist) / float64(totalBits)

	fmt.Printf("Исходное сообщение: %q\n", message)
	fmt.Printf("Изменённый байт #%d, бит #%d\n", byteIndex, bitIndex)
	fmt.Printf("Хеш 1: %x\n", hash1)
	fmt.Printf("Хеш 2: %x\n", hash2)
	fmt.Printf("Изменилось %d бит из %d (%.2f%%)\n", dist, totalBits, ratio*100)
}
