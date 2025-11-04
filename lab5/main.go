package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"os"
	"strings"

	"golang.org/x/image/bmp"
)

func messageToBits(msg string) []bool {
	data := []byte(msg)
	var bits []bool
	for _, b := range data {
		for i := 7; i >= 0; i-- {
			bits = append(bits, (b>>i)&1 == 1)
		}
	}
	return bits
}

func bitsToMessage(bits []bool) string {
	if len(bits)%8 != 0 {
		bits = bits[:len(bits)-(len(bits)%8)]
	}
	var msg strings.Builder
	for i := 0; i < len(bits); i += 8 {
		var b byte
		for j := 0; j < 8; j++ {
			if bits[i+j] {
				b |= 1 << (7 - j)
			}
		}
		msg.WriteByte(b)
	}
	return msg.String()
}

func embedLSB(img *image.RGBA, bits []bool) {
	if len(bits) == 0 {
		return
	}

	msgLenBits := uint32(len(bits))
	lenBits := make([]bool, 32)
	for i := 31; i >= 0; i-- {
		lenBits[31-i] = (msgLenBits>>i)&1 == 1
	}

	allBits := append(lenBits, bits...)

	idx := 0
	bounds := img.Bounds()
	for y := bounds.Max.Y - 1; y >= bounds.Min.Y; y-- {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if idx >= len(allBits) {
				return
			}
			r, g, b, a := img.At(x, y).RGBA()
			rr := uint8(r >> 8)
			gg := uint8(g >> 8)
			bb := uint8(b >> 8)
			aa := uint8(a >> 8)

			if allBits[idx] {
				rr |= 1
			} else {
				rr &^= 1
			}

			img.Set(x, y, color.RGBA{R: rr, G: gg, B: bb, A: aa})
			idx++
		}
	}
}

func extractLSB(img *image.RGBA) string {
	bounds := img.Bounds()
	var allBits []bool

	for y := bounds.Max.Y - 1; y >= bounds.Min.Y; y-- {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, _, _, _ := img.At(x, y).RGBA()
			rr := uint8(r >> 8)
			allBits = append(allBits, rr&1 == 1)
		}
	}

	if len(allBits) < 32 {
		return ""
	}

	var msgLen uint32
	for i := 0; i < 32; i++ {
		if allBits[i] {
			msgLen |= 1 << (31 - i)
		}
	}

	if int(msgLen) > len(allBits)-32 {
		msgLen = uint32(len(allBits) - 32)
	}

	msgBits := allBits[32 : 32+msgLen]
	return bitsToMessage(msgBits)
}

func loadBMPImage(filename string) (*image.RGBA, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	b := img.Bounds()
	rgba := image.NewRGBA(b)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			rgba.Set(x, y, img.At(x, y))
		}
	}
	return rgba, nil
}

func saveBMPImage(filename string, img image.Image) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	return bmp.Encode(file, img)
}

func main() {
	fmt.Print("Введите сообщение для внедрения: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	message := scanner.Text()

	img, err := loadBMPImage("input.bmp")
	if err != nil {
		fmt.Printf("Ошибка загрузки изображения: %v\n", err)
		return
	}

	bits := messageToBits(message)

	embedLSB(img, bits)

	err = saveBMPImage("output.bmp", img)
	if err != nil {
		fmt.Printf("Ошибка сохранения изображения: %v\n", err)
		return
	}

	fmt.Println("Сообщение успешно внедрено в output.bmp")

	img2, err := loadBMPImage("output.bmp")
	if err != nil {
		fmt.Printf("Ошибка загрузки для извлечения: %v\n", err)
		return
	}

	extracted := extractLSB(img2)
	fmt.Printf("Извлеченное сообщение: %s\n", extracted)
}
