package wav

import (
	"bufio"
	bin "encoding/binary"
	"os"
)

type WavData struct {
	bChunkID  [4]byte // B
	ChunkSize uint32  // L
	bFormat   [4]byte // B

	bSubchunk1ID  [4]byte // B
	Subchunk1Size uint32  // L
	AudioFormat   uint16  // L
	NumChannels   uint16  // L
	SampleRate    uint32  // L
	ByteRate      uint32  // L
	BlockAlign    uint16  // L
	BitsPerSample uint16  // L

	bSubchunk2ID  [4]byte // B
	Subchunk2Size uint32  // L
	Data          []byte  // L
}

func NewWavData(fn string) (*WavData, error) {
	res, err := os.OpenFile(fn, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	file := bufio.NewReader(res)

	wav := &WavData{}
	bin.Read(file, bin.BigEndian, &wav.bChunkID)
	bin.Read(file, bin.LittleEndian, &wav.ChunkSize)
	bin.Read(file, bin.BigEndian, &wav.bFormat)

	bin.Read(file, bin.BigEndian, &wav.bSubchunk1ID)
	bin.Read(file, bin.LittleEndian, &wav.Subchunk1Size)
	bin.Read(file, bin.LittleEndian, &wav.AudioFormat)
	bin.Read(file, bin.LittleEndian, &wav.NumChannels)
	bin.Read(file, bin.LittleEndian, &wav.SampleRate)
	bin.Read(file, bin.LittleEndian, &wav.ByteRate)
	bin.Read(file, bin.LittleEndian, &wav.BlockAlign)
	bin.Read(file, bin.LittleEndian, &wav.BitsPerSample)

	bin.Read(file, bin.BigEndian, &wav.bSubchunk2ID)
	bin.Read(file, bin.LittleEndian, &wav.Subchunk2Size)

	wav.Data = make([]byte, wav.Subchunk2Size)
	bin.Read(file, bin.LittleEndian, &wav.Data)

	return wav, nil
}

func (w *WavData) SampleCount() int {
	return int(len(w.Data) / 2)
}

func (w *WavData) Sample(index int) int16 {
	in := index * 2
	return btoi16(w.Data[in : in+2])
}

func btoi16(b []byte) int16 {
	value := int16(b[0])
	value += int16(b[1]) << 8
	return value
}
