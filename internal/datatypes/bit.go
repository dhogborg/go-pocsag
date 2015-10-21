package datatypes

// A bit is high or low, 0 or 1, true or false.
type Bit bool

// Int returns the bit value as 0 or 1
func (b Bit) Int() int {
	if b {
		return 1
	} else {
		return 0
	}
}

func (b Bit) UInt8() uint8 {
	return uint8(b.Int())
}
