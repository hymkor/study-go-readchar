//go:build run
// +build run

package main

// 準標準ライブラリだけで作った readline のデモコード
//
//   https://gist.github.com/artyom/a59e2707976124f387f5
// のコードを
// - Windows 対応
// - golang.org/x/crypto/ssh/terminal -> golang.org/x/term 化
// した。Windows で動作確認したところ
// - 矢印キーなど ASCII コードをもたないキーは認識しない
// - 漢字など2桁セル分の幅を占める文字があるとカーソル位置がずれる
// などの問題があった。 Emacs風のCtrlキー入力は大丈夫だった。
// - WSL では矢印キーはちゃんと認識された。

import (
	"fmt"
	"io"
	"log"
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
	screen := struct {
		io.Reader
		io.Writer
	}{os.Stdin, os.Stdout}
	term := term.NewTerminal(screen, "")
	term.SetPrompt(string(term.Escape.Red) + "> " + string(term.Escape.Reset))

	rePrefix := string(term.Escape.Cyan) + "Human says:" + string(term.Escape.Reset)

	for {
		line, err := term.ReadLine()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if line == "" {
			continue
		}
		fmt.Fprintln(term, rePrefix, line)
	}
}

func main() {
	if err := mains(); err != nil {
		log.Fatal(err)
	}
}
