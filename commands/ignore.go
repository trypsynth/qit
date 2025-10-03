package commands

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/trypsynth/qit/utils"
)

func NewIgnoreCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "ignore <template_name|list>",
		Short: "download .gitignore template(s) from gitignore.io or list available templates",
		RunE: func(cmd *cobra.Command, args []string) error {
			utils.RequireArgs(args, "Missing template name(s). Use 'ignore list' to see available templates.")
			if strings.ToLower(args[0]) == "list" {
				return listGitignoreTemplates()
			}
			return downloadGitignore(args[0])
		},
	}
}

func downloadGitignore(templates string) error {
	url := fmt.Sprintf("https://www.toptal.com/developers/gitignore/api/%s", templates)
	return utils.HTTPRequest(url, func(resp *http.Response) error {
		if resp.StatusCode == 200 {
			body, err := utils.ReadBody(resp)
			if err != nil {
				return err
			}
			if err := os.WriteFile(".gitignore", []byte(body), 0644); err != nil {
				return err
			}
			fmt.Printf("Downloaded .gitignore for %s\n", templates)
		} else {
			fmt.Fprintf(os.Stderr, "Error: HTTP %d from gitignore.io\n", resp.StatusCode)
			os.Exit(1)
		}
		return nil
	})
}

func listGitignoreTemplates() error {
	url := "https://www.toptal.com/developers/gitignore/api/list?format=lines"
	return utils.HTTPRequest(url, func(resp *http.Response) error {
		if resp.StatusCode == 200 {
			body, err := utils.ReadBody(resp)
			if err != nil {
				return err
			}
			fmt.Println("Available gitignore templates:")
			fmt.Print(body)
		} else {
			fmt.Fprintf(os.Stderr, "Error fetching list. HTTP %d\n", resp.StatusCode)
		}
		return nil
	})
}
