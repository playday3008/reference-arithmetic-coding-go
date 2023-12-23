package frequency_table

import (
	"fmt"
)

type SimpleFrequencyTable struct {
	frequencies []uint
	cumulative  []uint
	total       uint
}

var reinit_cumulative = false

func (self *SimpleFrequencyTable) checkSymbol(symbol uint) {
	if symbol >= uint(len(self.frequencies)) {
		panic("Symbol out of range")
	}
}

func (self *SimpleFrequencyTable) initCumulative(check_total ...any) { // Kind of a hack to allow optional arguments
	var sum uint = 0
	for i := range self.frequencies {
		sum += self.frequencies[i]
		self.cumulative[i+1] = sum
	}

	if len(check_total) <= 0 && sum != self.total {
		panic("Total frequency mismatch")
	}

	if reinit_cumulative {
		reinit_cumulative = false
	}
}

func (self *SimpleFrequencyTable) String() string {
	return fmt.Sprintf("SimpleFrequencyTable(frequencies=%v)", self.frequencies)
}

func NewSimpleFrequencyTable(freqs *FrequencyTable) *SimpleFrequencyTable {
	if freqs == nil {
		panic("Frequency table is null")
	}

	size := (*freqs).GetSymbolLimit()
	if size < 1 {
		panic("At least 1 symbol required")
	}

	self := &SimpleFrequencyTable{
		frequencies: make([]uint, size),
		cumulative:  make([]uint, size+1),
		total:       0,
	}
	for i := uint(0); i < size; i++ {
		self.frequencies[i] = (*freqs).Get(i)
	}
	self.initCumulative(true)
	self.total = (*freqs).GetHigh(size - 1)
	return self
}

func (self *SimpleFrequencyTable) GetSymbolLimit() uint {
	return uint(len(self.frequencies))
}

func (self *SimpleFrequencyTable) Get(symbol uint) uint {
	self.checkSymbol(symbol)
	return self.frequencies[symbol]
}

func (self *SimpleFrequencyTable) Set(symbol uint, freq uint) {
	self.checkSymbol(symbol)
	if freq < 0 {
		panic("Negative frequency")
	}

	temp := self.total - self.frequencies[symbol]
	if temp > self.total { // temp is uint and is always >= 0
		panic("Total underflow")
	}

	self.total = temp + freq
	self.frequencies[symbol] = freq
	reinit_cumulative = true
}

func (self *SimpleFrequencyTable) Increment(symbol uint) {
	self.checkSymbol(symbol)
	self.total++
	self.frequencies[symbol]++
	reinit_cumulative = true
}

func (self *SimpleFrequencyTable) GetTotal() uint {
	return self.total
}

func (self *SimpleFrequencyTable) GetLow(symbol uint) uint {
	self.checkSymbol(symbol)
	if reinit_cumulative {
		self.initCumulative()
	}
	return self.cumulative[symbol]
}

func (self *SimpleFrequencyTable) GetHigh(symbol uint) uint {
	self.checkSymbol(symbol)
	if reinit_cumulative {
		self.initCumulative()
	}
	return self.cumulative[symbol+1]
}
