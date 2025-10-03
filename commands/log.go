package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/trypsynth/qit/utils"
)

func NewLogCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "log",
		Short: "show commit history in readable format",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := utils.Git(false, "log", fmt.Sprintf("--pretty=format:%s", utils.CommitFormat), fmt.Sprintf("--date=format:%s", utils.DateFormat))
			if err == nil {
				fmt.Println()
			}
			return err
		},
	}
}
