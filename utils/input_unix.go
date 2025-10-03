//go:build !windows

package utils

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

func PromptForKey(prompt string) (rune, error) {
	fmt.Print(prompt)
	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return 0, fmt.Errorf("failed to set raw mode: %w", err)
	}
	defer term.Restore(fd, oldState)
	buf := make([]byte, 1)
	_, err = os.Stdin.Read(buf)
	if err != nil {
		return 0, fmt.Errorf("failed to read key: %w", err)
	}
	return rune(buf[0]), nil
}
