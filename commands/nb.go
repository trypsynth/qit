package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/trypsynth/qit/utils"
)

func NewNbCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "nb <branch_name>",
		Short: "switch to the branch with  the specified name, creating it if it doesn't exist",
		RunE: func(cmd *cobra.Command, args []string) error {
			utils.RequireArgs(args, "Missing branch name.")
			name := args[0]
			current, err := utils.CurrentBranch()
			if err != nil {
				return err
			}
			if current == name {
				fmt.Printf("Already on branch %s.\n", name)
				return nil
			}
			exists, err := utils.BranchExists(name)
			if err != nil {
				return err
			}
			if exists {
				fmt.Printf("Branch %s already exists. Switching to it...\n", name)
				return utils.Git(false, "checkout", name)
			}
			return utils.Git(false, "checkout", "-b", name)
		},
	}
}
