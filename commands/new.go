package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/trypsynth/qit/utils"
)

func NewNewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "new",
		Short: "pull and list recent commits",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			oldHead, err := utils.GitOutput("rev-parse", "HEAD")
			if err != nil {
				return err
			}
			if err := utils.Git(true, "pull"); err != nil {
				return err
			}
			newHead, err := utils.GitOutput("rev-parse", "HEAD")
			if err != nil {
				return err
			}
			if oldHead == newHead {
				fmt.Println("Nothing new.")
			} else {
				fmt.Println("Commits since last pull:")
				err := utils.Git(false, "log", fmt.Sprintf("%s..%s", oldHead, newHead), fmt.Sprintf("--pretty=format:%s", utils.CommitFormat), fmt.Sprintf("--date=format:%s", utils.DateFormat))
				if err == nil {
					fmt.Println()
				}
				return err
			}
			return nil
		},
	}
}
