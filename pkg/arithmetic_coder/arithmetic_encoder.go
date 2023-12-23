package arithmetic_coder

import (
	"github.com/spacemeshos/bitstream"

	"github.com/playday3008/reference-arithmetic-coding-go/pkg/frequency_table"
)

type ArithmeticEncoder struct {
	ArithmeticCoder
	output        *bitstream.BitWriter
	num_underflow uint
}

func NewArithmeticEncoder(num_bits int, output bitstream.BitWriter) *ArithmeticEncoder {
	self := &ArithmeticEncoder{}
	self.ArithmeticCoder = *NewArithmeticCoder(num_bits)
	self.output = &output
	self.num_underflow = 0
	self.SetVTable(self)
	return self
}

func (self *ArithmeticEncoder) Write(freqs *frequency_table.FrequencyTable, symbol uint) {
	if freqs == nil {
		panic("Frequency table is null")
	}

	self.Update(freqs, symbol)
}

func (self *ArithmeticEncoder) Finish() {
	self.output.WriteBit(bitstream.One)
	self.output.Flush(bitstream.Zero)
}

func (self *ArithmeticEncoder) Shift() {
	bit := self.low >> (self.num_state_bits - 1)
	if bit == 1 {
		self.output.WriteBit(bitstream.One)
	} else if bit == 0 {
		self.output.WriteBit(bitstream.Zero)
	} else {
		panic("Arithmetic underflow")
	}

	// Write out the saved underflow bits
	for ; self.num_underflow > 0; self.num_underflow-- {
		// Save as self.output.WriteBit(bit ^ 1)
		if bit^1 == 1 {
			self.output.WriteBit(bitstream.One)
		} else if bit^1 == 0 {
			self.output.WriteBit(bitstream.Zero)
		}
	}
}

func (self *ArithmeticEncoder) Underflow() {
	self.num_underflow++
}
