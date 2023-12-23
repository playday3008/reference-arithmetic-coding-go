package arithmetic_coder

import (
	"math"

	"github.com/playday3008/reference-arithmetic-coding-go/pkg/frequency_table"
)

type ArithmeticCoderVirtual interface {
	Shift()
	Underflow()
}

type ArithmeticCoder struct {
	num_state_bits int
	full_range     uint64
	half_range     uint64
	quarter_range  uint64
	minimum_range  uint64
	maximum_total  uint64
	state_mask     uint64

	low  uint64
	high uint64

	virtual ArithmeticCoderVirtual
}

func (self *ArithmeticCoder) SetVTable(v ArithmeticCoderVirtual) {
	self.virtual = v
}

// Constructs an arithmetic coder, which initializes the code range.
func NewArithmeticCoder(num_bits int) *ArithmeticCoder {
	if !(1 <= num_bits && num_bits <= 63) {
		panic("State size out of range")
	}
	self := &ArithmeticCoder{}
	self.num_state_bits = num_bits
	self.full_range = uint64(1) << self.num_state_bits
	self.half_range = self.full_range >> 1
	self.quarter_range = self.half_range >> 1
	self.minimum_range = self.quarter_range + 2
	self.maximum_total = min(math.MaxUint64/self.full_range, self.minimum_range)
	self.state_mask = self.full_range - 1

	self.low = 0
	self.high = self.state_mask
	return self
}

// Updates the code range (low and high) of this arithmetic coder as a result of processing the specified symbol with the specified frequency table.
func (self *ArithmeticCoder) Update(freqs *frequency_table.FrequencyTable, symbol uint) {
	// Null check
	if freqs == nil {
		panic("Frequency table is null")
	}

	// State check
	if self.low >= self.high || (self.low&self.state_mask) != self.low || (self.high&self.state_mask) != self.high {
		panic("Low or high out of range")
	}

	scope := self.high - self.low + 1
	if !(self.minimum_range <= scope && scope <= self.full_range) {
		panic("Scope out of range")
	}

	// Frequency table values check
	total := (*freqs).GetTotal()
	symbol_low := (*freqs).GetLow(symbol)
	symbol_high := (*freqs).GetHigh(symbol)
	if symbol_low == symbol_high {
		panic("Symbol has zero frequency")
	}
	if uint64(total) > self.maximum_total {
		panic("Cannot code symbol because total is too large")
	}

	// Update range
	new_low := self.low + uint64(symbol_low)*scope/uint64(total)
	new_high := self.low + uint64(symbol_high)*scope/uint64(total) - 1
	self.low = new_low
	self.high = new_high

	// While the highest bits are equal
	for ((self.low ^ self.high) & self.half_range) == 0 {
		self.virtual.Shift()
		self.low = (self.low << 1) & self.state_mask
		self.high = ((self.high << 1) & self.state_mask) | 1
	}

	// While the second highest bit of low is 1 and the second highest bit of high is 0
	for (self.low &^ self.high & self.quarter_range) != 0 {
		self.virtual.Underflow()
		self.low = (self.low << 1) ^ self.half_range
		self.high = ((self.high ^ self.half_range) << 1) | self.half_range | 1
	}
}

// Called to handle the situation when the top bit of 'low' and 'high' are equal.
func (self *ArithmeticCoder) Shift() {
	panic("Not implemented")
}

// Called to handle the situation when low=01(...) and high=10(...).
func (self *ArithmeticCoder) Underflow() {
	panic("Not implemented")
}
