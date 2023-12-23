package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"runtime/pprof"
	"time"

	"github.com/spacemeshos/bitstream"

	. "github.com/playday3008/reference-arithmetic-coding-go/pkg/arithmetic_coder"
	. "github.com/playday3008/reference-arithmetic-coding-go/pkg/frequency_table"
)

var action string
var input_file string
var output_file string
var cpu_profile string
var show_stats bool

const (
	// Number of bits for the arithmetic coding range. Must be in the range [1, 62].
	num_bits = 32
	// Maximum frequency table capacity. Must be at least 2, and at most 2^32 - 1.
	num_symbols = 257
)

const (
	EOF = 256
)

func init() {
	var overwrite bool
	flag.StringVar(&action, "action", "", "Action to perform: compress, decompress")
	flag.StringVar(&input_file, "input", "", "Input file")
	flag.StringVar(&output_file, "output", "", "Output file")
	flag.BoolVar(&overwrite, "overwrite", false, "Overwrite output file if it exists")
	flag.BoolVar(&show_stats, "show_stats", false, "Show statistics")
	flag.StringVar(&cpu_profile, "cpu_profile", "", "Write cpu profile to file")
	flag.Parse()

	if *flag.Bool("help", false, "Print help") == true {
		flag.PrintDefaults()
		os.Exit(0)
	}

	if action == "" {
		fmt.Println("Action not specified")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if action != "compress" && action != "decompress" {
		fmt.Println("Invalid action: ", action)
		flag.PrintDefaults()
		os.Exit(2)
	}

	if input_file == "" {
		fmt.Println("Input file not specified")
		flag.PrintDefaults()
		os.Exit(3)
	}

	if output_file == "" {
		fmt.Println("Output file not specified")
		flag.PrintDefaults()
		os.Exit(4)
	}

	if _, err := os.Stat(input_file); os.IsNotExist(err) {
		fmt.Println("Input file does not exist")
		flag.PrintDefaults()
		os.Exit(5)
	}

	if _, err := os.Stat(output_file); !os.IsNotExist(err) && !overwrite {
		fmt.Println("Output file already exists")
		var answer string
		fmt.Printf("Overwrite %s? (y/n) ", output_file)
		fmt.Scanf("%s", &answer)
		if answer != "y" {
			os.Exit(6)
		}
	}
}

func main() {
	fmt.Println("Action: ", action)
	fmt.Println("Input file: ", input_file)
	fmt.Println("Output file: ", output_file)

	if cpu_profile != "" {
		f, err := os.Create(cpu_profile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if show_stats {
		start_time := time.Now()
		defer func() {
			fmt.Println("\tTime elapsed: ", time.Since(start_time))
		}()
	}

	if action == "compress" {
		compress(input_file, output_file)
	} else if action == "decompress" {
		decompress(input_file, output_file)
	}

	fmt.Println("Done")

	if show_stats {
		print_stats(input_file, output_file, action)
	}
}

func compress(input_file string, output_file string) {
	input, err := os.Open(input_file)
	if err != nil {
		log.Fatal(err)
	}
	defer input.Close()

	output, err := os.Create(output_file)
	if err != nil {
		log.Fatal(err)
	}
	defer output.Close()

	bufWriter := bufio.NewWriter(output)
	defer bufWriter.Flush()
	bitWriter := bitstream.NewWriter(bufWriter)
	// bitWriter.Flush() done by Finish()

	var flat_freqs FrequencyTable = NewFlatFrequencyTable(num_symbols)
	var freqs FrequencyTable = NewSimpleFrequencyTable(&flat_freqs)
	enc := NewArithmeticEncoder(num_bits, *bitWriter)
	defer enc.Finish()
	defer enc.Write(&freqs, EOF)

	b := make([]byte, 1)
	for {
		n, err := input.Read(b)
		if err != nil || n != 1 {
			break
		}
		enc.Write(&freqs, uint(b[0]))
		freqs.Increment(uint(b[0]))
	}
}

func decompress(input_file string, output_file string) {
	input, err := os.Open(input_file)
	if err != nil {
		log.Fatal(err)
	}
	defer input.Close()

	output, err := os.Create(output_file)
	if err != nil {
		log.Fatal(err)
	}
	defer output.Close()

	bufReader := bufio.NewReader(input)
	bitReader := bitstream.NewReader(bufReader)

	var flat_freqs FrequencyTable = NewFlatFrequencyTable(num_symbols)
	var freqs FrequencyTable = NewSimpleFrequencyTable(&flat_freqs)
	dec := NewArithmeticDecoder(num_bits, *bitReader)
	for {
		symbol := dec.Read(&freqs)
		if symbol == EOF {
			break
		}
		output.Write([]byte{byte(symbol)})
		freqs.Increment(symbol)
	}
}

func print_stats(input_file string, output_file string, action string) {
	entropy_input, err := file_entropy(input_file)
	if err != nil {
		log.Fatal(err)
	}
	entropy_output, err := file_entropy(output_file)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Statistics: ")
	fmt.Println("\tEntropy:")
	fmt.Println("\t\tInput entropy:\t", entropy_input)
	fmt.Println("\t\tOutput entropy:\t", entropy_output)
	if action == "compress" {
		fmt.Println("\t\tEntropy ratio:\t", entropy_output/entropy_input)
	} else if action == "decompress" {
		fmt.Println("\t\tEntropy ratio:\t", entropy_input/entropy_output)
	}

	input_stat, err := os.Stat(input_file)
	if err != nil {
		log.Fatal(err)
	}
	output_stat, err := os.Stat(output_file)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\tFile size:")
	fmt.Println("\t\tInput file size:\t", input_stat.Size())
	fmt.Println("\t\tOutput file size:\t", output_stat.Size())

	if action == "compress" {
		fmt.Println("\tCompression ratio: ", float64(input_stat.Size())/float64(output_stat.Size()))
	} else if action == "decompress" {
		fmt.Println("\t\tFile size ratio:\t", float64(output_stat.Size())/float64(input_stat.Size()))
	}
}

func file_entropy(name string) (float64, error) {
	file, err := os.ReadFile(name)
	if err != nil {
		return 0, err
	}

	byte_map := make(map[byte]int)
	for _, b := range file {
		byte_map[b]++
	}

	var entropy float64 = 0
	for _, v := range byte_map {
		p := float64(v) / float64(len(file))
		if p > 0 {
			entropy -= p * math.Log2(p)
		}
	}

	return entropy, nil
}
