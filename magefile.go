//go:build mage

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/magefile/mage/sh"
)

var Default = Build

func Build() error {
	binary := "qit"
	if runtime.GOOS == "windows" {
		binary = "qit.exe"
	}
	return sh.Run("go", "build", "-o", binary, ".")
}

func BuildAll() error {
	platforms := []struct {
		os   string
		arch string
		ext  string
	}{
		{"windows", "amd64", ".exe"},
		{"linux", "amd64", ""},
		{"darwin", "amd64", ""},
		{"darwin", "arm64", ""},
	}
	for _, p := range platforms {
		binary := fmt.Sprintf("qit-%s-%s%s", p.os, p.arch, p.ext)
		env := map[string]string{
			"GOOS":   p.os,
			"GOARCH": p.arch,
		}
		if err := sh.RunWith(env, "go", "build", "-o", binary, "."); err != nil {
			return err
		}
	}
	return nil
}

func Clean() error {
	patterns := []string{
		"qit",
		"qit.exe",
		"qit-*",
	}
	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return err
		}
		for _, match := range matches {
			if err := os.Remove(match); err != nil && !os.IsNotExist(err) {
				return err
			}
		}
	}
	return nil
}

func Install() error {
	return sh.Run("go", "install", ".")
}

func Fmt() error {
	return sh.Run("go", "fmt", "./...")
}

func Vet() error {
	return sh.Run("go", "vet", "./...")
}
