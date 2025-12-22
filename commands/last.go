package commands

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/trypsynth/qit/utils"
)

func NewLastCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "last [<number>]",
		Short: "show the last <number> commits. Defaults to 1",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			count := 1
			if len(args) > 0 {
				if n, err := strconv.Atoi(args[0]); err == nil {
					count = n
				}
			}
			err := utils.Git(false, "log", fmt.Sprintf("-%d", count), fmt.Sprintf("--pretty=format:%s", utils.CommitFormat), fmt.Sprintf("--date=format:%s", utils.DateFormat))
			if err == nil {
				fmt.Println()
			}
			return err
		},
	}
}
