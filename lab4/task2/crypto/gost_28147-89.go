package crypto

import (
	"crypto/rand"
	"encoding/binary"
	"math/bits"
)

const (
	KEY_LENGTH         = 32
	KEY_BLOCKS         = 8
	MESSAGE_BLOCK      = 8
	CIPHER_ROUNDS      = 24
	CIPHER_LAST_ROUNDS = 8
)

var sboxes = [8][16]byte{
	{4, 10, 9, 2, 13, 8, 0, 14, 6, 11, 1, 12, 7, 15, 5, 3},
	{14, 11, 4, 12, 6, 13, 15, 10, 2, 3, 8, 1, 0, 7, 5, 9},
	{5, 8, 1, 13, 10, 3, 4, 2, 14, 15, 12, 7, 6, 0, 9, 11},
	{7, 13, 10, 1, 0, 8, 9, 15, 14, 4, 6, 12, 11, 2, 5, 3},
	{6, 12, 7, 1, 5, 15, 13, 8, 4, 10, 9, 14, 0, 3, 11, 2},
	{4, 11, 10, 0, 7, 2, 1, 13, 3, 6, 8, 5, 9, 12, 15, 14},
	{13, 11, 4, 1, 3, 15, 5, 9, 0, 10, 14, 7, 6, 8, 2, 12},
	{1, 15, 13, 0, 5, 7, 10, 4, 9, 2, 3, 14, 6, 11, 8, 12},
}

type Blocks struct {
	blocks []*Block
}

func NewBlocks(message []byte) *Blocks {
	messageLength := len(message)
	blocksCount := calculateBlocks(messageLength)
	blocks := make([]*Block, blocksCount)

	for i := 0; i < blocksCount; i++ {
		blocks[i] = NewBlock(message[i*MESSAGE_BLOCK:])
	}

	return &Blocks{blocks: blocks}
}

func (blocks *Blocks) GetData() []byte {
	var data []byte
	for _, block := range blocks.blocks {
		data = append(data, block.ToBytes()...)
	}

	return data
}

type Block struct {
	HasExtra bool
	L        [4]byte
	R        [4]byte
}

func NewBlock(message []byte) *Block {
	hasExtra := false
	var extraByte byte = 0
	messageLength := len(message)

	if messageLength < MESSAGE_BLOCK {
		hasExtra = true
		extraByte = byte(MESSAGE_BLOCK - messageLength)
	}

	var A []byte
	if messageLength < MESSAGE_BLOCK {
		A = append(A, message[0:messageLength]...)
		for i := messageLength; i < 8; i++ {
			A = append(A, extraByte)
		}
	} else {
		A = append(A, message[0:MESSAGE_BLOCK]...)
	}

	return &Block{HasExtra: hasExtra, L: [4]byte(A[0:4]), R: [4]byte(A[4:8])}
}

func (block *Block) ToBytes() []byte {
	var message []byte

	var count int = 0
	if block.HasExtra {
		extra := block.R[3]
		count = int(extra)
	}

	message = append(message, block.L[:]...)
	message = append(message, block.R[:]...)

	return message[0 : MESSAGE_BLOCK-count]
}

type Key struct {
	key [KEY_LENGTH]byte
}

func (key *Key) GetBlock(i int) [4]byte {
	blockNum := i % KEY_BLOCKS
	var block [4]byte
	copy(block[:], key.key[blockNum*4:blockNum*4+4])

	return block
}

func GenerateKey() (*Key, error) {
	var key [KEY_LENGTH]byte

	_, err := rand.Read(key[:])
	if err != nil {
		return nil, err
	}

	return &Key{key}, nil
}

func calculateBlocks(messageLength int) int {
	blocksCount := messageLength / MESSAGE_BLOCK
	if messageLength%MESSAGE_BLOCK > 0 {
		blocksCount += 1
	}

	return blocksCount
}

func Encrypt(message []byte, key *Key) []byte {
	blocks := NewBlocks(message)

	for _, block := range blocks.blocks {
		EncryptBlock(block, key)
	}

	return blocks.GetData()
}

func EncryptBlock(block *Block, key *Key) {
	for i := 0; i < CIPHER_ROUNDS; i++ {
		TransformBlock(block, key.GetBlock(i))
	}

	for i := CIPHER_LAST_ROUNDS - 1; i >= 0; i-- {
		TransformBlock(block, key.GetBlock(i))
	}
}

func TransformBlock(block *Block, keyBlock [4]byte) {
	var carry uint32
	smod, _ := bits.Add32(bytesToUint32(block.R), bytesToUint32(keyBlock), carry)
	smodb := uint32ToBytes(smod)

	var ssimple [4]byte
	for j := 0; j < 8; j++ {
		sblock := (smodb[j/2] >> (4 * j)) & byte(15)
		v := sboxes[j][sblock]
		ssimple[j/2] |= v << (4 * j)
	}

	srol := bits.RotateLeft32(bytesToUint32(ssimple), 11)
	sxor := srol ^ bytesToUint32(block.L)

	block.L = block.R
	block.R = uint32ToBytes(sxor)
}

func DecryptBlock(block *Block, key *Key) {
	for i := 0; i < CIPHER_LAST_ROUNDS; i++ {
		TransformBlock(block, key.GetBlock(i))
	}

	for i := CIPHER_ROUNDS - 1; i >= 0; i-- {
		TransformBlock(block, key.GetBlock(i))
	}
}

func Decrypt(encrypted []byte, key *Key) []byte {
	blocks := NewBlocks(encrypted)

	for _, block := range blocks.blocks {
		DecryptBlock(block, key)
	}

	return blocks.GetData()
}

func bytesToUint32(b [4]byte) uint32 {
	return binary.LittleEndian.Uint32(b[:])
}

func uint32ToBytes(u uint32) [4]byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, u)

	return [4]byte(b)
}
