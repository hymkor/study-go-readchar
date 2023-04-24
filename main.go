package main

import (
	"bufio"
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
	disable, err := enable(stdin)
	if err != nil {
		return err
	}
	defer disable()

	defer term.Restore(stdin, oldState)

	w := bufio.NewWriter(os.Stdout)
	for {
		var buffer [256]byte

		n, err := os.Stdin.Read(buffer[:])
		if err != nil {
			return err
		}
		ch := buffer[:n]
		for _, c := range string(ch) {
			if c <= 32 {
				fmt.Fprintf(w, "\\x%02X ", c)
			} else {
				fmt.Fprintf(w, "%c", c)
			}
		}
		w.Flush()
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
