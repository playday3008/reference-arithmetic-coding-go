package frequency_table

import (
	"fmt"
)

type FlatFrequencyTable struct {
	num_symbols uint
}

func (self *FlatFrequencyTable) checkSymbol(symbol uint) {
	if symbol >= self.num_symbols {
		panic("Symbol out of range")
	}
}

func (self *FlatFrequencyTable) String() string {
	return fmt.Sprintf("FlatFrequencyTable(numSymbols=%d)", self.num_symbols)
}

func NewFlatFrequencyTable(numSymbols uint) *FlatFrequencyTable {
	if numSymbols < 1 {
		panic("Number of symbols must be positive")
	}
	return &FlatFrequencyTable{numSymbols}
}

func (self *FlatFrequencyTable) GetSymbolLimit() uint {
	return self.num_symbols
}

func (self *FlatFrequencyTable) Get(symbol uint) uint {
	self.checkSymbol(symbol)
	return 1
}

func (self *FlatFrequencyTable) Set(symbol uint, freq uint) {
	panic("Not implemented")
}

func (self *FlatFrequencyTable) Increment(symbol uint) {
	panic("Not implemented")
}

func (self *FlatFrequencyTable) GetTotal() uint {
	return self.num_symbols
}

func (self *FlatFrequencyTable) GetLow(symbol uint) uint {
	self.checkSymbol(symbol)
	return symbol
}

func (self *FlatFrequencyTable) GetHigh(symbol uint) uint {
	self.checkSymbol(symbol)
	return symbol + 1
}
