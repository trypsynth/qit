//go:build !windows
// +build !windows

package utils

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

func PromptForKey(prompt string) rune {
	fmt.Print(prompt)
	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return '\x00'
	}
	defer term.Restore(fd, oldState)
	buf := make([]byte, 1)
	_, err = os.Stdin.Read(buf)
	if err != nil {
		return '\x00'
	}
	return rune(buf[0])
}
