//go:build windows
// +build windows

package utils

import (
	"fmt"
	"syscall"
)

var (
	msvcrt    = syscall.NewLazyDLL("msvcrt.dll")
	procGetch = msvcrt.NewProc("_getch")
)

func PromptForKey(prompt string) rune {
	fmt.Print(prompt)
	ret, _, _ := procGetch.Call()
	return rune(ret)
}
