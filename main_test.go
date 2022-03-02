package main

import (
	"bytes"
	"math/rand"
	"os"
	"path"
	"strconv"
	"testing"
)

func TestToStdout(t *testing.T) {
	args := []string{"tee"}

	want := "hello world"
	stdin := bytes.NewBuffer([]byte(want))

	var stdout bytes.Buffer

	err := run(args, stdin, &stdout)
	if err != nil {
		t.Errorf("Unexpected error: %q", err)
	}

	got := stdout.String()
	if got != want {
		t.Errorf("Wrong output: got %q, want %q", got, want)
	}
}

func TestToFile(t *testing.T) {
	temp := os.TempDir()
	filenames := []string{
		path.Join(temp, "tee"+strconv.Itoa(rand.Int())+".txt"),
		path.Join(temp, "tee"+strconv.Itoa(rand.Int())+".txt"),
	}
	defer func() {
		for _, filename := range filenames {
			_ = os.Remove(filename)
		}
	}()

	args := []string{"tee"}
	args = append(args, filenames...)

	want := "hello world"
	stdin := bytes.NewBuffer([]byte(want))

	var stdout bytes.Buffer
	if err := run(args, stdin, &stdout); err != nil {
		t.Errorf("Unexpected error: %q", err)
	}

	for _, filename := range filenames {
		bytes, err := os.ReadFile(filename)
		if err != nil {
			t.Fatalf("failed to read temporary file %q", filename)
		}

		got := string(bytes)
		if got != want {
			t.Errorf("Wrong output: got %q, want %q", got, want)
		}
	}
}

func TestAppendToFile(t *testing.T) {
	temp := os.TempDir()
	filenames := []string{
		path.Join(temp, "tee"+strconv.Itoa(rand.Int())+".txt"),
		path.Join(temp, "tee"+strconv.Itoa(rand.Int())+".txt"),
	}
	defer func() {
		for _, filename := range filenames {
			_ = os.Remove(filename)
		}
	}()

	want := "hello world"

	for _, filename := range filenames {
		if err := os.WriteFile(filename, []byte(want[:3]), 0666); err != nil {
			t.Fatalf("failed to write to temporary file %q", filename)
		}
	}

	args := []string{"tee", "-a"}
	args = append(args, filenames...)

	stdin := bytes.NewBuffer([]byte(want[3:]))

	var stdout bytes.Buffer
	if err := run(args, stdin, &stdout); err != nil {
		t.Errorf("Unexpected error: %q", err)
	}

	for _, filename := range filenames {
		bytes, err := os.ReadFile(filename)
		if err != nil {
			t.Fatalf("failed to read temporary file %q", filename)
		}

		got := string(bytes)
		if got != want {
			t.Errorf("Wrong output: got %q, want %q", got, want)
		}
	}
}
