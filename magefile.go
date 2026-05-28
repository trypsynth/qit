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
	if err := os.MkdirAll("dist", 0755); err != nil {
		return err
	}
	platforms := []struct {
		os   string
		arch string
		ext  string
	}{
		{"windows", "amd64", ".exe"},
		{"windows", "arm64", ".exe"},
		{"linux", "amd64", ""},
		{"linux", "arm64", ""},
		{"darwin", "amd64", ""},
		{"darwin", "arm64", ""},
	}
	for _, p := range platforms {
		binary := filepath.Join("dist", fmt.Sprintf("qit-%s-%s%s", p.os, p.arch, p.ext))
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
	for _, name := range []string{"qit", "qit.exe"} {
		if err := os.Remove(name); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	return os.RemoveAll("dist")
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
