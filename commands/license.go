package commands

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/trypsynth/qit/utils"
)

func NewLicenseCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "license <license_name|list>",
		Short: "download a license template from GitHub or list available licenses",
		RunE: func(cmd *cobra.Command, args []string) error {
			utils.RequireArgs(args, "Missing license name or 'list'. Use 'license list' to see available licenses.")
			if strings.ToLower(args[0]) == "list" {
				return listGitHubLicenses()
			}
			return downloadGitHubLicense(strings.ToLower(args[0]))
		},
	}
}

func listGitHubLicenses() error {
	url := "https://api.github.com/licenses"
	return utils.HTTPRequest(url, func(resp *http.Response) error {
		if resp.StatusCode == 200 {
			var licenses []map[string]interface{}
			body, err := utils.ReadBody(resp)
			if err != nil {
				return err
			}
			if err := json.Unmarshal([]byte(body), &licenses); err != nil {
				return err
			}
			fmt.Println("Available licenses:")
			for _, license := range licenses {
				key := license["key"].(string)
				name := license["name"].(string)
				fmt.Printf("%s: %s\n", key, name)
			}
		} else {
			fmt.Fprintf(os.Stderr, "Error fetching licenses. HTTP %d\n", resp.StatusCode)
		}
		return nil
	})
}

func downloadGitHubLicense(licenseKey string) error {
	url := fmt.Sprintf("https://api.github.com/licenses/%s", licenseKey)
	return utils.HTTPRequest(url, func(resp *http.Response) error {
		if resp.StatusCode == 200 {
			body, err := utils.ReadBody(resp)
			if err != nil {
				return err
			}
			var licenseData map[string]interface{}
			if err := json.Unmarshal([]byte(body), &licenseData); err != nil {
				return err
			}
			licenseBody := licenseData["body"].(string)
			licenseName := licenseData["name"].(string)
			if err := os.WriteFile("LICENSE", []byte(licenseBody), 0644); err != nil {
				return err
			}
			fmt.Printf("Downloaded the %s license to LICENSE file.\n", licenseName)
		} else if resp.StatusCode == 404 {
			fmt.Fprintf(os.Stderr, "License '%s' not found. Use 'qit license list' to see available licenses.\n", licenseKey)
			os.Exit(1)
		} else {
			fmt.Fprintf(os.Stderr, "Error: HTTP %d from GitHub API\n", resp.StatusCode)
			os.Exit(1)
		}
		return nil
	})
}
