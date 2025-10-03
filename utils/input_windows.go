//go:build windows

package utils

import (
	"fmt"

	"golang.org/x/sys/windows"
)

var (
	msvcrt    = windows.NewLazySystemDLL("msvcrt.dll")
	procGetch = msvcrt.NewProc("_getch")
)

func PromptForKey(prompt string) (rune, error) {
	fmt.Print(prompt)
	ret, _, err := procGetch.Call()
	if ret == 0 {
		return 0, fmt.Errorf("failed to read key: %w", err)
	}
	return rune(ret), nil
}
