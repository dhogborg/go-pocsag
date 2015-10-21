package utils

import (
	. "gopkg.in/check.v1"
	"testing"

	"github.com/dhogborg/go-pocsag/internal/datatypes"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

var _ = Suite(&UtilitiesSuite{})

type UtilitiesSuite struct{}

func (f *UtilitiesSuite) Test_MSB_BitsToBytes_8_FF(c *C) {
	bits := []datatypes.Bit{
		true, true, true, true, true, true, true, true,
	}
	c.Assert(MSBBitsToBytes(bits, 8), DeepEquals, []byte{0xFF})
}

func (f *UtilitiesSuite) Test_MSB_BitsToBytes_16_FFFF(c *C) {
	bits := []datatypes.Bit{
		true, true, true, true, true, true, true, true,
		true, true, true, true, true, true, true, true,
	}
	c.Assert(MSBBitsToBytes(bits, 8), DeepEquals, []byte{0xFF, 0xFF})
}

func (f *UtilitiesSuite) Test_MSB_BitsToBytes_16_FF00(c *C) {
	bits := []datatypes.Bit{
		true, true, true, true, true, true, true, true,
		false, false, false, false, false, false, false, false,
	}
	c.Assert(MSBBitsToBytes(bits, 8), DeepEquals, []byte{0xFF, 0x00})
}

func (f *UtilitiesSuite) Test_LSB_BitsToBytes_8_FF(c *C) {
	bits := []datatypes.Bit{true, true, true, true, true, true, true, true}
	c.Assert(LSBBitsToBytes(bits, 8), DeepEquals, []byte{0xFF})
}

func (f *UtilitiesSuite) Test_LSB_BitsToBytes_16_FFFF(c *C) {
	bits := []datatypes.Bit{
		true, true, true, true, true, true, true, true,
		true, true, true, true, true, true, true, true,
	}
	c.Assert(LSBBitsToBytes(bits, 8), DeepEquals, []byte{0xFF, 0xFF})
}

func (f *UtilitiesSuite) Test_LSB_BitsToBytes_16_FF00(c *C) {
	bits := []datatypes.Bit{
		false, false, false, false, false, false, false, false,
		true, true, true, true, true, true, true, true,
	}
	c.Assert(LSBBitsToBytes(bits, 8), DeepEquals, []byte{0x00, 0xFF})
}

func (f *UtilitiesSuite) Test_LSB_BitsToBytes_16_00FF(c *C) {
	bits := []datatypes.Bit{
		true, true, true, true, true, true, true, true,
		false, false, false, false, false, false, false, false,
	}
	c.Assert(LSBBitsToBytes(bits, 8), DeepEquals, []byte{0xFF, 0x00})
}

// POCSAG speceific

func (f *UtilitiesSuite) Test_MSBBitsToBytes_Keyword_PREAMBLE(c *C) {
	bits := []datatypes.Bit{
		false, true, true, true,
		true, true, false, false,

		true, true, false, true,
		false, false, true, false,

		false, false, false, true,
		false, true, false, true,

		true, true, false, true,
		true, false, false, false,
	}
	c.Assert(MSBBitsToBytes(bits, 8), DeepEquals, []byte{0x7C, 0xD2, 0x15, 0xD8})
}

func (f *UtilitiesSuite) Test_MSBBitsToBytes_Keyword_IDLE(c *C) {
	bits := []datatypes.Bit{
		false, true, true, true,
		true, false, true, false,

		true, false, false, false,
		true, false, false, true,

		true, true, false, false,
		false, false, false, true,

		true, false, false, true,
		false, true, true, true,
	}
	c.Assert(MSBBitsToBytes(bits, 8), DeepEquals, []byte{0x7A, 0x89, 0xC1, 0x97})
}

func (f *UtilitiesSuite) Test_BCD_Min(c *C) {
	bits := []datatypes.Bit{
		false, false, false, false,
	}
	c.Assert(BitcodedDecimals(bits), Equals, "0")
}

func (f *UtilitiesSuite) Test_BCD_Max(c *C) {
	bits := []datatypes.Bit{
		true, true, true, true,
	}
	c.Assert(BitcodedDecimals(bits), Equals, "(")
}

func (f *UtilitiesSuite) Test_BCD_10chars(c *C) {

	bits := []datatypes.Bit{
		false, false, false, false,
		true, true, true, false,
		false, false, false, false,
		true, true, true, false,

		true, false, false, false,
		true, false, false, true,

		true, true, false, false,
		true, true, false, false,

		false, false, false, true,
		true, false, true, false,
	}

	c.Assert(BitcodedDecimals(bits), Equals, "0707193385")
}
