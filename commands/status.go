package commands

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

func NewStatusCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "show simplified summary of working directory changes",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			staged, unstaged, err := getGitStatus()
			if err != nil {
				return err
			}
			if len(staged) == 0 && len(unstaged) == 0 {
				fmt.Println("Working tree clean.")
				return nil
			}
			if len(staged) > 0 {
				fmt.Println("Staged for commit:")
				for _, f := range staged {
					fmt.Printf("  %s\n", f)
				}
			}
			if len(unstaged) > 0 {
				fmt.Println("Not staged for commit:")
				for _, f := range unstaged {
					fmt.Printf("  %s\n", f)
				}
			}
			return nil
		},
	}
}

func getGitStatus() ([]string, []string, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to get git status.")
		os.Exit(1)
	}
	var staged []string
	var unstaged []string
	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}
		if len(line) < 3 {
			continue
		}
		x := line[0]
		y := line[1]
		file := strings.TrimSpace(line[3:])
		if x != ' ' && x != '?' {
			staged = append(staged, file)
		}
		if y != ' ' {
			unstaged = append(unstaged, file)
		}
	}
	return staged, unstaged, scanner.Err()
}
