package utils

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	CommitFormat = "%h %an: %s (%ad)."
	DateFormat   = "%Y-%m-%d %H:%M:%S"
	UserAgent   = "qit-cli"
)

func Git(quiet bool, args ...string) error {
	cmd := exec.Command("git", args...)
	if !quiet {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git %s failed: %w", strings.Join(args, " "), err)
	}
	return nil
}

func RequireArgs(args []string, missingMessage string) error {
	if len(args) == 0 {
		return errors.New(missingMessage)
	}
	return nil
}

func HTTPGet(url string) (*http.Response, error) {
	return HTTPGetWithHeaders(url, nil)
}

func HTTPGetWithHeaders(url string, headers map[string]string) (*http.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30 * time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch from %s: %w", url, err)
	}
	return resp, nil
}

func CurrentBranch() (string, error) {
	return GitOutput("rev-parse", "--abbrev-ref", "HEAD")
}

func BranchExists(name string) (bool, error) {
	output, err := GitOutput("branch", "--list", name)
	if err != nil {
		return false, err
	}
	return len(output) > 0, nil
}

func GitOutput(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git %s failed: %w", strings.Join(args, " "), err)
	}
	return strings.TrimSpace(string(output)), nil
}

func ReadBody(resp *http.Response) ([]byte, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	return body, nil
}
