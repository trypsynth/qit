package commands

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/trypsynth/qit/utils"
)

func NewCpCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "cp <message>",
		Short: "commit changes to tracked files with the specified message, and push",
		RunE: func(cmd *cobra.Command, args []string) error {
			utils.RequireArgs(args, "Missing commit message.")
			message := strings.Join(args, " ")
			if err := utils.Git(false, "commit", "-am", message); err != nil {
				return err
			}
			return utils.Git(false, "push")
		},
	}
}
