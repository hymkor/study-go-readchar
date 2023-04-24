package main

import (
	"bytes"
	"fmt"
	"os"

	"golang.org/x/term"
)

func mains() error {
	stdin := int(os.Stdin.Fd())
	stdout := int(os.Stdout.Fd())
	if !term.IsTerminal(stdin) || !term.IsTerminal(stdout) {
		return fmt.Errorf("stdin/stdout should be terminal")
	}
	oldState, err := term.MakeRaw(stdin)
	if err != nil {
		return err
	}
	defer term.Restore(stdin, oldState)
	for {
		var buffer [256]byte

		n, err := os.Stdin.Read(buffer[:])
		if err != nil {
			return err
		}
		ch := buffer[:n]
		fmt.Print(ch)
		// exit with Ctrl-Z
		if bytes.Equal(ch, []byte{'z' & 0x1F}) {
			return nil
		}
	}
}

func main() {
	if err := mains(); err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}
}
