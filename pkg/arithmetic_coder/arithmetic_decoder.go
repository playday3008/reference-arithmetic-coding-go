package arithmetic_coder

import (
	"github.com/spacemeshos/bitstream"

	"github.com/playday3008/reference-arithmetic-coding-go/pkg/frequency_table"
)

type ArithmeticDecoder struct {
	ArithmeticCoder
	input *bitstream.BitReader
	code  uint64
}

func (self *ArithmeticDecoder) readCodeBit() uint {
	bit, err := self.input.ReadBit()
	if err != nil || bit == bitstream.Zero {
		return 0
	}
	return 1
}

func NewArithmeticDecoder(num_bits int, input bitstream.BitReader) *ArithmeticDecoder {
	self := &ArithmeticDecoder{}
	self.ArithmeticCoder = *NewArithmeticCoder(num_bits)
	self.input = &input
	self.code = 0
	for i := 0; i < self.num_state_bits; i++ {
		self.code = (self.code << 1) | uint64(self.readCodeBit())
	}
	self.SetVTable(self)
	return self
}

func (self *ArithmeticDecoder) Read(freqs *frequency_table.FrequencyTable) uint {
	if freqs == nil {
		panic("Frequency table is null")
	}

	total := (*freqs).GetTotal()
	if uint64(total) > self.maximum_total {
		panic("Cannot decode symbol because total is too large")
	}

	// Translate from coding range scale to frequency table scale
	scope := self.high - self.low + 1
	offset := self.code - self.low
	value := ((offset+1)*uint64(total) - 1) / scope
	if value*scope/uint64(total) > offset {
		panic("Assertion error: Overflow")
	} else if value >= uint64(total) {
		panic("Assertion error: Value too high")
	}

	// A kind of binary search. Find highest symbol such that freqs.cumulative[symbol] <= offset.
	// This works because the frequency table entries are monotonically non-decreasing.
	var start uint = 0
	var end uint = (*freqs).GetSymbolLimit()
	for end-start > 1 {
		middle := (start + end) >> 1
		if uint64((*freqs).GetLow(middle)) > value {
			end = middle
		} else {
			start = middle
		}
	}
	if start+1 != end {
		panic("Assertion error: Expected start + 1 == end")
	}

	symbol := start
	if !(uint64((*freqs).GetLow(symbol))*scope/uint64(total) <= offset && offset < uint64((*freqs).GetHigh(symbol))*scope/uint64(total)) {
		panic("Assertion error: Symbol out of range")
	}
	self.Update(freqs, symbol)
	if !(self.low <= self.code && self.code <= self.high) {
		panic("Assertion error: Code out of range")
	}
	return symbol
}

func (self *ArithmeticDecoder) Shift() {
	self.code = ((self.code << 1) & self.state_mask) | uint64(self.readCodeBit())
}

func (self *ArithmeticDecoder) Underflow() {
	self.code = (self.code & self.half_range) | ((self.code << 1) & (self.state_mask >> 1)) | uint64(self.readCodeBit())
}
