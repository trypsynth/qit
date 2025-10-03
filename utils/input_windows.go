//go:build windows

package utils

import (
	"errors"
	"fmt"

	"golang.org/x/sys/windows"
)

var (
	msvcrt    = windows.NewLazySystemDLL("msvcrt.dll")
	procGetch = msvcrt.NewProc("_getch")
)

func PromptForKey(prompt string) (rune, error) {
	fmt.Print(prompt)
	ret, _, _ := procGetch.Call()
	if ret == 0 {
		return 0, errors.New("failed to read key")
	}
	return rune(ret), nil
}
