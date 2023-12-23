package frequency_table

type FrequencyTable interface {
	String() string
	// Returns the number of symbols in this frequency table.
	GetSymbolLimit() uint
	// Returns the frequency of the given symbol.
	Get(symbol uint) uint
	// Sets the frequency of the given symbol to the given value.
	Set(symbol uint, freq uint)
	// Increments the frequency of the given symbol.
	Increment(symbol uint)
	// Returns the total of all symbol frequencies.
	// The returned value is always equal to getHigh(getSymbolLimit() - 1).
	GetTotal() uint
	// Returns the sum of the frequencies of all the symbols strictly below the given symbol value.
	GetLow(symbol uint) uint
	// Returns the sum of the frequencies of the given symbol and all the symbols below.
	GetHigh(symbol uint) uint
}
