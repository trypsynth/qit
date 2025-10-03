package commands

import (
	"fmt"
	"os"
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
			confirm := utils.PromptForKey("This will discard all changes. Continue? (y/n) ")
			fmt.Println()
			if strings.ToLower(string(confirm)) != "y" {
				os.Exit(0)
			}
			return utils.Git(false, "reset", "--hard")
		},
	}
}
