package lsb

import (
	"encoding/binary"
	"image"
	"image/color"
)

// Hide writes data into copy of src using LSB (least significant bit) method.
// It includes 32 bit length value in data before payload
func Hide(src image.Image, data []byte) image.Image {
	bounds := src.Bounds()
	out := image.NewRGBA(bounds)
	var currentBit uint32
	dataLen := uint32(len(data))
	allData := make([]byte, 4)
	binary.LittleEndian.PutUint32(allData, dataLen)
	allData = append(allData, data...)
	dataLen = uint32(len(allData))

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := src.At(x, y).RGBA()
			values := toUint8Array(r, g, b, a)
			currentByte := currentBit / 8
			// Leave alpha channel untouched
			for i := 0; i < 3 && currentByte < dataLen; i++ {
				values[i] = setBit(values[i], bitAt(allData[currentByte], currentBit%8), 0)
				currentBit++
				currentByte = currentBit / 8
			}
			newColor := color.RGBA{
				R: values[0],
				G: values[1],
				B: values[2],
				A: values[3],
			}
			out.Set(x, y, newColor)
		}
	}
	return out
}

// Reveal reads length header and payload from src.
// Only payload is included in returned data
func Reveal(src image.Image) []byte {
	data := make([]byte, 0)
	var dataLen uint32 = 4
	bounds := src.Bounds()
	var currentBit uint32
	var buffer byte
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := src.At(x, y).RGBA()
			colorChannels := toUint8Array(r, g, b, a)
			currentByte := currentBit / 8
			for i := 0; i < 3 && currentByte < dataLen; i++ {
				buffer = setBit(buffer, bitAt(colorChannels[i], 0), uint8(currentBit%8))
				currentBit++
				// One whole byte has been written to buffer
				if currentBit%8 == 0 {
					data = append(data, buffer)
					// Four fist bytes are read.
					// These contain the length of the payload in bytes
					if len(data) == 4 {
						dataLen = binary.LittleEndian.Uint32(data)
						currentBit = 0
					}
				}
				currentByte = currentBit / 8
			}
		}
	}
	return data[4:]
}

func toUint8Array(values ...uint32) (result []uint8) {
	result = make([]uint8, len(values))
	for i, value := range values {
		result[i] = uint8(value >> 8)
	}
	return
}

func bitAt(val byte, pos uint32) uint8 {
	return uint8((val >> pos) & 1)
}
func setBit(target uint8, val uint8, pos uint8) uint8 {
	target = (target & ^uint8(1<<pos)) | (val << pos)
	return target
}
