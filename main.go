package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

const (
	// exitFailure is the exit code if the program terminates
	// on a failure.
	exitFailure = 1
)

func openAll(filenames []string, add bool) ([]*os.File, error) {
	perm := os.O_WRONLY | os.O_CREATE
	if add {
		perm |= os.O_APPEND
	} else {
		perm |= os.O_TRUNC
	}

	var files []*os.File
	for _, filename := range filenames {
		file, err := os.OpenFile(filename, perm, 0666)
		if err != nil {
			closeAll(files)
			return nil, err
		}
		files = append(files, file)
	}

	return files, nil
}

func closeAll(files []*os.File) {
	for _, file := range files {
		_ = file.Close()
	}
}

func makeWriterSlice(w io.Writer, files ...*os.File) []io.Writer {
	r := make([]io.Writer, 0, len(files))
	r = append(r, w)

	for _, f := range files {
		r = append(r, f)
	}
	return r
}

func run(args []string, stdin io.Reader, stdout io.Writer) error {
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	flags.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: tee [options] [files ...]\n")
		fmt.Fprintf(os.Stderr, "options:\n")
		flags.PrintDefaults()
	}

	add := flags.Bool("a", false, "Appends the output to each file, instead of overwriting.")
	// @TODO: handle interrupts
	//ignore := flags.Bool("i", false, "Ignores interrupts.")
	if err := flags.Parse(args[1:]); err != nil {
		return err
	}

	files, err := openAll(flags.Args(), *add)
	if err != nil {
		return err
	}
	defer closeAll(files)

	w := io.MultiWriter(makeWriterSlice(stdout, files...)...)
	r := io.TeeReader(stdin, w)
	_, err = io.ReadAll(r)

	return err
}

func main() {
	if err := run(os.Args, os.Stdin, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(exitFailure)
	}
}
