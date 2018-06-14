package main

import (
	"bufio"
	"fmt"
	"log"
	"strings"

	"github.com/alexflint/go-arg"
	"github.com/brentp/xopen"
)

// Opts is the struct with the options that the program accepts.
type Opts struct {
	Input string `arg:"positional,required" help:"input file/s (STDIN if -)"`
	RNA   bool   `arg:"" help:"convert to RNA"`
}

// Version returns the program version.
func (Opts) Version() string { return "fastq-rnadna 0.1" }

// Description returns an extended description of the program.
func (Opts) Description() string {
	return "Convert sequence representation of a FASTQ file to DNA or RNA."
}

func main() {
	var opts Opts
	arg.MustParse(&opts)

	// open input file.
	f, err := xopen.Ropen(opts.Input)
	if err != nil {
		log.Fatal(err)
	}

	// create scanner.
	sc := bufio.NewScanner(f)

	idx := 0 // determines line index within a fastq block.
	for sc.Scan() {
		l := sc.Text()
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

	if err = sc.Err(); err != nil {
		log.Fatal(sc.Err())
	}

	if err = f.Close(); err != nil {
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
