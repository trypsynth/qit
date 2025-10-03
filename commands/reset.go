package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/trypsynth/qit/utils"
)

func NewResetCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "reset",
		Short: "hard reset to last commit, discarding all changes",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			confirm, err := utils.PromptForKey("This will discard all changes. Continue? (y/n) ")
			fmt.Println()
			if err != nil {
				return err
			}
			if strings.ToLower(string(confirm)) != "y" {
				return fmt.Errorf("operation cancelled")
			}
			return utils.Git(false, "reset", "--hard")
		},
	}
}
