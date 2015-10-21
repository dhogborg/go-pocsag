package pocsag

import (
	"bytes"
	"fmt"

	"github.com/dhogborg/go-pocsag/internal/wav"
)

// ReadWav reads a wav file from disc and puts it in memory for the
// scanner to parse as a standard transmission
func ReadWav(path string) *bytes.Buffer {

	wavdata, err := wav.NewWavData(path)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	samplecount := int(wavdata.Subchunk2Size / uint32(wavdata.BitsPerSample/8))
	seconds := float32(samplecount) / float32(wavdata.SampleRate)

	if DEBUG {
		fmt.Printf("Samples: %d\n", samplecount)
		fmt.Printf("Seconds: %0.3f\n", seconds)
	}

	buffer := bytes.NewBuffer(wavdata.Data)
	return buffer

}
