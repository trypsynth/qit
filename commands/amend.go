package commands

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/trypsynth/qit/utils"
)

func NewAmendCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "amend <message>",
		Short: "amend the last commit with a new message",
		RunE: func(cmd *cobra.Command, args []string) error {
			utils.RequireArgs(args, "Missing new commit message.")
			message := strings.Join(args, " ")
			return utils.Git(false, "commit", "--amend", "--reset", "-m", message)
		},
	}
}
