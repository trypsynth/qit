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

type gitHubLicense struct {
	Key  string `json:"key"`
	Name string `json:"name"`
	Body string `json:"body"`
}

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
		"User-Agent": utils.UserAgent,
	}
	resp, err := utils.HTTPGetWithHeaders(url, headers)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error fetching licenses: HTTP %d", resp.StatusCode)
	}
	body, err := utils.ReadBody(resp)
	if err != nil {
		return err
	}
	var licenses []gitHubLicense
	if err := json.Unmarshal(body, &licenses); err != nil {
		return fmt.Errorf("failed to parse licenses: %w", err)
	}
	fmt.Println("Available licenses:")
	for _, license := range licenses {
		fmt.Printf("%s: %s\n", license.Key, license.Name)
	}
	return nil
}

func downloadGitHubLicense(licenseKey string) error {
	url := fmt.Sprintf("https://api.github.com/licenses/%s", licenseKey)
	headers := map[string]string{
		"User-Agent": utils.UserAgent,
	}
	resp, err := utils.HTTPGetWithHeaders(url, headers)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("license '%s' not found, use 'qit license list' to see available licenses", licenseKey)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d from GitHub API", resp.StatusCode)
	}
	body, err := utils.ReadBody(resp)
	if err != nil {
		return err
	}
	var license gitHubLicense
	if err := json.Unmarshal(body, &license); err != nil {
		return fmt.Errorf("failed to parse license: %w", err)
	}
	if license.Body == "" {
		return fmt.Errorf("invalid license data: missing body field")
	}
	if err := os.WriteFile("LICENSE", []byte(license.Body), 0644); err != nil {
		return err
	}
	fmt.Printf("Downloaded the %s license to LICENSE file.\n", license.Name)
	return nil
}
