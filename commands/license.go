package commands

import (
	"encoding/json"
	"fmt"
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
			if err := utils.RequireArgs(args, "missing license name or 'list', use 'license list' to see available licenses"); err != nil {
				return err
			}
			if strings.ToLower(args[0]) == "list" {
				return listGitHubLicenses()
			}
			return downloadGitHubLicense(strings.ToLower(args[0]))
		},
	}
}

func listGitHubLicenses() error {
	url := "https://api.github.com/licenses"
	headers := map[string]string{
		"User-Agent": "qit-cli",
	}
	resp, err := utils.HTTPGetWithHeaders(url, headers)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("error fetching licenses: HTTP %d", resp.StatusCode)
	}
	var licenses []map[string]interface{}
	body, err := utils.ReadBody(resp)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, &licenses); err != nil {
		return err
	}
	fmt.Println("Available licenses:")
	for _, license := range licenses {
		key, ok := license["key"].(string)
		if !ok {
			continue
		}
		name, ok := license["name"].(string)
		if !ok {
			continue
		}
		fmt.Printf("%s: %s\n", key, name)
	}
	return nil
}

func downloadGitHubLicense(licenseKey string) error {
	url := fmt.Sprintf("https://api.github.com/licenses/%s", licenseKey)
	headers := map[string]string{
		"User-Agent": "qit-cli",
	}
	resp, err := utils.HTTPGetWithHeaders(url, headers)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 404 {
		return fmt.Errorf("license '%s' not found, use 'qit license list' to see available licenses", licenseKey)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("HTTP %d from GitHub API", resp.StatusCode)
	}
	body, err := utils.ReadBody(resp)
	if err != nil {
		return err
	}
	var licenseData map[string]interface{}
	if err := json.Unmarshal(body, &licenseData); err != nil {
		return err
	}
	licenseBody, ok := licenseData["body"].(string)
	if !ok {
		return fmt.Errorf("invalid license data: missing or invalid 'body' field")
	}
	licenseName, ok := licenseData["name"].(string)
	if !ok {
		return fmt.Errorf("invalid license data: missing or invalid 'name' field")
	}
	if err := os.WriteFile("LICENSE", []byte(licenseBody), 0644); err != nil {
		return err
	}
	fmt.Printf("Downloaded the %s license to LICENSE file.\n", licenseName)
	return nil
}
