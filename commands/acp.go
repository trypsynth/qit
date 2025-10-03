package commands

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/trypsynth/qit/utils"
)

func NewAcpCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "acp <message>",
		Short: "add all files, commit with the specified message, and push",
		RunE: func(cmd *cobra.Command, args []string) error {
			utils.RequireArgs(args, "Missing commit message")
			message := strings.Join(args, " ")
			if err := utils.Git(false, "add", "."); err != nil {
				return err
			}
			if err := utils.Git(false, "commit", "-m", message); err != nil {
				return err
			}
			return utils.Git(false, "push")
		},
	}
}
