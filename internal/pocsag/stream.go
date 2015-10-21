package pocsag

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/dhogborg/go-pocsag/internal/datatypes"
	"github.com/dhogborg/go-pocsag/internal/utils"
)

type StreamReader struct {
	Stream *bufio.Reader
	// 0 for auto
	baud int
}

// NewStreamReader returns a new stream reader for the source provided.
// Set bauds 0 for automatic detection.
func NewStreamReader(source io.Reader, bauds int) *StreamReader {

	return &StreamReader{
		Stream: bufio.NewReader(source),
		baud:   bauds,
	}

}

// StartScan takes a channel on which bitstreams will be written when found and parsed.
// The scanner will continue indefently and sleep for 3 ms per cycle to go easy on the system load.
func (s *StreamReader) StartScan(bitstream chan []datatypes.Bit) {

	fmt.Println("Starting transmission scanner")

	for {

		bytes := make([]byte, 4096)
		c, err := s.Stream.Read(bytes)

		if err != nil {
			println(err.Error())
			os.Exit(0)
		}

		stream := s.bToInt16(bytes[:c])

		start, bitlength := s.ScanTransmissionStart(stream)

		if start > 0 {

			blue.Println("-- Transmission received at", time.Now(), "--------------")

			transmission := s.ReadTransmission(stream[start:])

			bits := utils.StreamToBits(transmission, bitlength)
			bitstream <- bits
		}

		time.Sleep(3 * time.Millisecond)

	}
}

// ReadTransmission reads the beginning and subsequent datapackages into
// a new buffer until encounters noise instead of signal.
func (s *StreamReader) ReadTransmission(beginning []int16) []int16 {

	stream := make([]int16, 0)
	stream = append(stream, beginning...)

	for {

		bytes := make([]byte, 4096)
		c, _ := s.Stream.Read(bytes)

		if c > 0 {

			bstr := s.bToInt16(bytes[:c])
			stream = append(stream, bstr...)

			if s.isNoise(bstr) {
				break
			}
		}

	}

	return stream
}

// ScanTransmissionStart scans for repeated 1010101010101 pattern of bits in the
// stream. A minimum of 400 samples is required to sucessfully sync the receiver with
// the stream. ScanTransmissionStart looks for when the signal wave changes from high to low
// and reversed, and measures the distance between those changes. They should correspond to
// the bitlength determined by the current baud-rate. An attempt at guessing the baudrate is
// also made when a repeated pattern is found.
// retuned is the index of the stream on which the caller should begin reading bits, and
// the estimated bitlength, the number of samples between each bit center in transmission stream.
func (s *StreamReader) ScanTransmissionStart(stream []int16) (int, int) {

	if len(stream) == 0 {
		return -1, 0
	}

	switches := []int{}
	prevsamp := stream[0]
	first_switch := -1

	// find the indexes where we cross the 0-boundary
	// if we switch sides we store the index in an array for further analasys
	for a, sample := range stream {

		if (prevsamp > 0 && sample < 0) || (prevsamp < 0 && sample > 0) {
			switches = append(switches, a)
			if first_switch < 0 {
				first_switch = a
			}
		}

		prevsamp = sample
	}

	// find the mean distance between boundary corsses
	sum := 0.0
	for a := 0; a < len(switches)-1; a += 1 {
		sum += float64(switches[a+1] - switches[a])
	}

	mean_bitlength := sum / float64(len(switches)-1)

	bitlength := float64(s.bitlength(int(mean_bitlength)))

	// if bitlength is not on a scale of known baudrates then
	// we probably don't have a pocsag sync-transmission
	if bitlength < 0 {
		return -1, 0
	}

	if DEBUG {
		fmt.Println("Mean bitlength:", mean_bitlength)
		fmt.Println("Determined bitlength:", bitlength)
	}

	// look at every other sample to see if we have a repeating pattern with matching size
	confidence := 0
	for a := 0; a < len(switches)-3; a += 1 {

		// length from switch a to a+1
		w1 := float64(switches[a+1] - switches[a])
		w2 := float64(switches[a+3] - switches[a+2])

		// how much the persumed bits vary from eachother
		intravariance := (w1 / w2) - 1
		if intravariance < 0 {
			intravariance = intravariance * -1
		}

		// how much the persumed bits vary from the determined bitlength
		baudvariance := (w1 / bitlength) - 1
		if baudvariance < 0 {
			baudvariance = baudvariance * -1
		}

		// don't stray more than 20%
		if intravariance < 0.2 && baudvariance < 0.2 {
			confidence += 1
		} else {
			confidence = 0
		}

		if confidence > 10 {

			if DEBUG {
				fmt.Println("Found bitsync")
			}

			return switches[a] + int(bitlength/2), int(bitlength)
		}

	}

	return -1, 0
}

// bitlength returns the proper bitlength from a calcualated mean distance between
// wave transitions. If the baudrate is set by configuration then that is used instead.
func (s *StreamReader) bitlength(mean int) int {

	if mean > 150 && mean < 170 {
		return 160
	} else if mean > 75 && mean < 85 || s.baud == 600 {
		return 80
	} else if mean > 35 && mean < 45 || s.baud == 1200 {
		return 40
	} else if mean > 15 && mean < 25 || s.baud == 2400 {
		return 20
	} else {
		return -1
	}

}

// isNoise detects noise by calculating the number of times the signal goes over the 0-line
// during a signal this value is between 25 and 50, but noise is above 100, usually around 300-400.
func (s *StreamReader) isNoise(stream []int16) bool {

	if len(stream) == 0 {
		return false
	}

	prevsamp := stream[0]
	switches := 0

	// find the indexes where we cross the 0-boundary
	for _, sample := range stream {

		if (prevsamp > 0 && sample < 0) || (prevsamp < 0 && sample > 0) {
			switches += 1
		}

		prevsamp = sample
	}

	return switches > 100
}

// bToInt16 converts bytes to int16
func (s *StreamReader) bToInt16(b []byte) (u []int16) {
	u = make([]int16, len(b)/2)
	for i, _ := range u {
		val := int16(b[i*2])
		val += int16(b[i*2+1]) << 8
		u[i] = val
	}
	return
}
