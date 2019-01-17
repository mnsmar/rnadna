package main

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/alexflint/go-arg"
	"github.com/brentp/xopen"
)

// Opts is the struct with the options that the program accepts.
type Opts struct {
	Input string `arg:"positional,required" help:"input file (STDIN if -)"`
	RNA   bool   `arg:"" help:"convert to RNA"`
}

// Version returns the program version.
func (Opts) Version() string { return "fastq-rnadna 0.2" }

// Description returns an extended description of the program.
func (Opts) Description() string {
	return "Convert sequence representation of a FASTQ file to DNA or RNA."
}

func main() {
	var opts Opts
	arg.MustParse(&opts)

	// open input file.
	r, err := xopen.Ropen(opts.Input)
	if err != nil {
		log.Fatal(err)
	}

	idx := 0 // determines line index within a fastq block.
	for {
		l, err := readLn(r)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}

		if idx%4 == 1 {
			if opts.RNA {
				l = strings.Map(d2r, l)
			} else {
				l = strings.Map(r2d, l)
			}
		}
		fmt.Println(l)
		idx++
	}

	if err = r.Close(); err != nil {
		log.Fatal(err)
	}
}

func r2d(r rune) rune {
	switch r {
	case 'U':
		return 'T'
	case 'u':
		return 't'
	}
	return r
}

func d2r(r rune) rune {
	switch r {
	case 'T':
		return 'U'
	case 't':
		return 'u'
	}
	return r
}

// readLn returns a single line (without the ending \n) from the input reader.
// An error is returned if an error is returned by the reader.
func readLn(r *xopen.Reader) (string, error) {
	var (
		isPrefix bool  = true
		err      error = nil
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln), err
}
