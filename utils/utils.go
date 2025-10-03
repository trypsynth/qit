package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

const (
	CommitFormat = "%h %an: %s (%ad)."
	DateFormat   = "%Y-%m-%d %H:%M:%S"
)

func Git(quiet bool, args ...string) error {
	cmd := exec.Command("git", args...)
	if !quiet {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git %s failed", strings.Join(args, " "))
	}
	return nil
}

func ErrorExit(message string) {
	fmt.Fprintln(os.Stderr, message)
	os.Exit(1)
}

func RequireArgs(args []string, missingMessage string) {
	if len(args) == 0 {
		ErrorExit(missingMessage)
	}
}

func HTTPRequest(url string, callback func(*http.Response) error) error {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to fetch from %s: %v\n", url, err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	return callback(resp)
}

func CurrentBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func BranchExists(name string) (bool, error) {
	cmd := exec.Command("git", "branch", "--list", name)
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}
	return len(strings.TrimSpace(string(output))) > 0, nil
}

func GitOutput(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func ReadBody(resp *http.Response) (string, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
