package pocsag

import (
	. "gopkg.in/check.v1"
	"testing"

	"github.com/dhogborg/go-pocsag/internal/datatypes"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

var _ = Suite(&PocsagSuite{})

type PocsagSuite struct{}

func (f *PocsagSuite) Test_Syndrome_Calculation_Normal_Address(c *C) {
	// valid bitstream
	bits := bitstream("01010001111011110011110111000010")
	syndr := syndrome(bits)

	c.Assert(syndr > 0, Equals, false)
}

func (f *PocsagSuite) Test_Syndrome_Calculation_Normal_Idle(c *C) {
	// valid bitstream
	bits := bitstream("01111010100010011100000110010111")
	syndr := syndrome(bits)

	c.Assert(syndr > 0, Equals, false)
}

func (f *PocsagSuite) Test_Syndrome_Calculation_Normal_Message(c *C) {
	// valid bitstream
	bits := bitstream("11001101100000000000011110001100")
	syndr := syndrome(bits)

	c.Assert(syndr > 0, Equals, false)
}

func (f *PocsagSuite) Test_Syndrome_Calculation_Error_1(c *C) {
	//                      v error
	bits := bitstream("01010101111011110011110111000010")
	syndr := syndrome(bits)

	c.Assert(syndr > 0, Equals, true)
}

func (f *PocsagSuite) Test_Syndrome_Calculation_Error_2(c *C) {
	//                   v  v errors
	bits := bitstream("01110101111011110011110111000010")
	syndr := syndrome(bits)

	c.Assert(syndr > 0, Equals, true)

}

func (f *PocsagSuite) Test_BitCorrection_No_Rrror(c *C) {

	bits := bitstream("01010001111011110011110111000010")
	cbits, corr := BitCorrection(bits)
	stream := streambits(cbits)

	c.Assert(corr, Equals, 0)
	c.Assert(stream, Equals, "01010001111011110011110111000010")
}

func (f *PocsagSuite) Test_BitCorrection_PayloadError(c *C) {
	//                      v error
	bits := bitstream("01010101111011110011110111000010")
	cbits, corr := BitCorrection(bits)
	stream := streambits(cbits)

	c.Assert(corr, Equals, 1)
	c.Assert(stream, Equals, "01010001111011110011110111000010")
}

func (f *PocsagSuite) Test_BitCorrection_PayloadErrors(c *C) {

	//                      v errors       v
	bits := bitstream("01010101111011110011010111000010")
	cbits, corr := BitCorrection(bits)
	stream := streambits(cbits)

	c.Assert(corr, Equals, 2)
	c.Assert(stream, Equals, "01010001111011110011110111000010")
}

func (f *PocsagSuite) Test_BitCorrection_ParityError(c *C) {
	//                                            v
	bits := bitstream("01010001111011110011110111010010")
	cbits, corr := BitCorrection(bits)
	stream := streambits(cbits)

	c.Assert(corr, Equals, 1)
	c.Assert(stream, Equals, "01010001111011110011110111000010")
}

func bitstream(stream string) []datatypes.Bit {
	bits := make([]datatypes.Bit, 32)
	for i, c := range stream {
		if string(c) == "1" {
			bits[i] = datatypes.Bit(true)
		} else {
			bits[i] = datatypes.Bit(false)
		}
	}
	return bits
}

func streambits(bits []datatypes.Bit) string {
	stream := ""
	for _, c := range bits {
		if c {
			stream += "1"
		} else {
			stream += "0"
		}
	}
	return stream
}
