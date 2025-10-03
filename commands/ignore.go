package commands

import (
	"fmt"
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
			if err := utils.RequireArgs(args, "missing template name(s), use 'ignore list' to see available templates"); err != nil {
				return err
			}
			if strings.ToLower(args[0]) == "list" {
				return listGitignoreTemplates()
			}
			return downloadGitignore(args[0])
		},
	}
}

func downloadGitignore(templates string) error {
	url := fmt.Sprintf("https://www.toptal.com/developers/gitignore/api/%s", templates)
	resp, err := utils.HTTPGet(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("HTTP %d from gitignore.io", resp.StatusCode)
	}
	body, err := utils.ReadBody(resp)
	if err != nil {
		return err
	}
	if err := os.WriteFile(".gitignore", body, 0644); err != nil {
		return err
	}
	fmt.Printf("Downloaded .gitignore for %s\n", templates)
	return nil
}

func listGitignoreTemplates() error {
	url := "https://www.toptal.com/developers/gitignore/api/list?format=lines"
	resp, err := utils.HTTPGet(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("error fetching list: HTTP %d", resp.StatusCode)
	}
	body, err := utils.ReadBody(resp)
	if err != nil {
		return err
	}
	fmt.Println("Available gitignore templates:")
	fmt.Print(string(body))
	return nil
}
