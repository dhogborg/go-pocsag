package datatypes

import (
	. "gopkg.in/check.v1"
	"testing"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

var _ = Suite(&TypeSuite{})

type TypeSuite struct{}

func (f *TypeSuite) Test_BitToInt(c *C) {
	high := Bit(true)
	low := Bit(false)
	c.Assert(high.Int(), Equals, 1)
	c.Assert(low.Int(), Equals, 0)
}

func (f *TypeSuite) Test_BitToUInt8(c *C) {
	high := Bit(true)
	low := Bit(false)
	c.Assert(high.UInt8(), Equals, uint8(1))
	c.Assert(low.UInt8(), Equals, uint8(0))
}
