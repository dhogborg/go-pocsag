package utils

import (
	"fmt"

	"github.com/fatih/color"

	"github.com/dhogborg/go-pocsag/internal/datatypes"
)

var (
	DEBUG bool
	LEVEL int
)

var (
	green = color.New(color.FgGreen)
	red   = color.New(color.FgRed)
	blue  = color.New(color.FgBlue)
)

// Tell the package to print debug data
func SetDebug(d bool, verbosity int) {
	DEBUG = d
	LEVEL = verbosity
}

// StreamToBits converts samples to bits using the bitlength specified.
// Observe that POCSAG signifies a high bit with a low frequency.
func StreamToBits(stream []int16, bitlength int) []datatypes.Bit {

	bits := make([]datatypes.Bit, (len(stream)/bitlength)+1)
	b := 0

	for a := 0; a < len(stream); a += bitlength {

		sample := stream[a]
		if a > 2 && a < len(stream)-2 {
			// let the samples before and after influence our sample, to prevent spike errors
			sample = (stream[a-1] / 2) + stream[a] + (stream[a+1] / 2)
		}

		bits[b] = datatypes.Bit((sample < 0))
		b += 1

	}

	return bits
}

// MSBBitsToBytes converts bitsream to bytes using MSB to LSB order.
func MSBBitsToBytes(bits []datatypes.Bit, bitsPerByte int) []byte {

	var b uint8
	bytes := []byte{}
	power := bitsPerByte - 1

	for a := 0; a < len(bits); a += 1 {

		bit := bits[a].UInt8()
		mod := a % bitsPerByte

		if mod == 0 && a > 0 {
			bytes = append(bytes, b)
			b = 0
		}

		pow := uint(power - mod)
		b += (bit * (1 << pow))

	}

	if len(bits)%bitsPerByte == 0 {
		bytes = append(bytes, b)
	}

	return bytes
}

// LSBBitsToBytes converts bitsream to bytes using LSB to MSB order.
func LSBBitsToBytes(bits []datatypes.Bit, bitsPerByte int) []byte {

	var b uint8
	bytes := []byte{}

	for a := 0; a < len(bits); a += 1 {

		bit := bits[a].UInt8()
		mod := a % bitsPerByte

		if mod == 0 && a > 0 {
			bytes = append(bytes, b)
			b = 0
		}

		pow := uint(mod)
		b += (bit * (1 << pow))

	}

	if len(bits)%bitsPerByte == 0 {
		bytes = append(bytes, b)
	}

	return bytes
}

// simple parity check
func ParityCheck(bits []datatypes.Bit, even_bit datatypes.Bit) bool {

	sum := even_bit.Int()
	for _, b := range bits {
		if b {
			sum += 1
		}
	}
	return (sum % 2) == 0
}

// BitcodedDecimals takes 4 bits per decimal to create values between 0 and 15.
// *) values 0-9 are used as is
// *) values 10-14 are special characters translated by bcdSpecial()
// *) value = 15 is not used.
func BitcodedDecimals(bits []datatypes.Bit) string {

	msg := ""
	var foo uint8 = 0

	bitsPerByte := 4

	for a := 0; a < len(bits); a += 1 {

		bit := bits[a].UInt8()
		mod := a % bitsPerByte

		if mod == 0 && a > 0 {
			msg += bcdChar(foo)
			foo = 0
		}

		pow := uint(mod)
		foo += (bit * (1 << pow))

	}

	if len(bits)%bitsPerByte == 0 {
		msg += bcdChar(foo)
	}

	return msg
}

// bcdChar translates digits and non-digit bitcoded entitis to charaters as per POCSAG protocol
func bcdChar(foo uint8) string {

	if foo < 10 {
		return fmt.Sprintf("%d", foo)
	}

	if foo == 10 {
		return ""
	}

	chars := []string{
		"",
		"U",
		" ",
		"-",
		")",
		"(",
	}
	return chars[foo-10]
}

func Btouint32(bytes []byte) uint32 {

	var a uint32 = 0
	a += uint32(bytes[0]) << 24
	a += uint32(bytes[1]) << 16
	a += uint32(bytes[2]) << 8
	a += uint32(bytes[3])

	return a
}

func TernaryStr(condition bool, a, b string) string {
	if condition {
		return a
	} else {
		return b
	}
}

// PrintStream, used for debugging of streams
func PrintStream(samples []int16) {
	for _, sample := range samples {
		PrintSample(sample)
	}
	println("")
}

// PrintBitstream, used for debugging of streams
func PrintBitstream(bits []datatypes.Bit) {
	for _, b := range bits {
		PrintSample(int16(b.Int()))
	}
	println("")
}

// PrintSample, used for debugging of streams
func PrintSample(sample int16) {
	if sample > 0 {
		green.Printf("%d ", sample)
	} else {
		red.Printf("%d ", sample)
	}
}

func PrintUint32(i uint32) {
	var x uint32 = 1 << 31
	for a := 0; a < 32; a += 1 {
		if (i & x) > 0 {
			green.Print("1 ")
		} else {
			red.Print("0 ")
		}
		x >>= 1
	}
	println("")
}
