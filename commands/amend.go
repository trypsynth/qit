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
			if err := utils.RequireArgs(args, "missing new commit message"); err != nil {
				return err
			}
			message := strings.Join(args, " ")
			return utils.Git(false, "commit", "--amend", "--reset", "-m", message)
		},
	}
}
