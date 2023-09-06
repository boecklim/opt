package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"opt/pkg/opt"
	"os"
	"path/filepath"
)

// Version is the command version, injected at build time.
var Version string = "dev"

type userFlags struct {
	outFile string
	remove  bool
	args    []string
}

func main() {
	var flags userFlags
	flag.StringVar(&flags.outFile, "out", "", "output file (default stdout)")
	// flag.StringVar(&flags.formatter, "fmt", "", "go pretty-printer: gofmt, goimports or noop (default gofmt)")
	printVersion := flag.Bool("version", false, "show the version for moq")
	flag.BoolVar(&flags.remove, "rm", false, "first remove output file, if it exists")

	flag.Usage = func() {
		fmt.Println(`opt [flags] source-dir struct-name`)
		flag.PrintDefaults()
	}

	flag.Parse()
	flags.args = flag.Args()

	if *printVersion {
		fmt.Printf("opt version %s\n", Version)
		os.Exit(0)
	}

	err := run(flags)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		flag.Usage()
		os.Exit(1)
	}
}

func run(flags userFlags) error {
	if len(flags.args) < 2 {
		return errors.New("not enough arguments")
	}

	if flags.remove && flags.outFile != "" {
		if err := os.Remove(flags.outFile); err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				return err
			}
		}
	}
	var buf bytes.Buffer
	var out io.Writer = os.Stdout
	if flags.outFile != "" {
		out = &buf
	}

	srcDir, args := flags.args[0], flags.args[1:]
	m, err := opt.New(opt.Config{
		SrcDir: srcDir,
	})
	if err != nil {
		return err
	}

	if err = m.Generate(out, args[0]); err != nil {
		return err
	}

	if flags.outFile == "" {
		return nil
	}

	// create the file
	err = os.MkdirAll(filepath.Dir(flags.outFile), 0o750)
	if err != nil {
		return err
	}

	return os.WriteFile(flags.outFile, buf.Bytes(), 0o600)
}
