package commands

import (
	"github.com/spf13/cobra"
	"github.com/trypsynth/qit/utils"
)

func NewUndoCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "undo",
		Short: "undo last commit while keeping changes intact",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return utils.Git(false, "reset", "--soft", "HEAD~1")
		},
	}
}
