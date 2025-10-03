package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/trypsynth/qit/utils"
)

func NewDbCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "db <branch_name>",
		Short: "delete the specified local branch",
		RunE: func(cmd *cobra.Command, args []string) error {
			utils.RequireArgs(args, "Missing branch name.")
			name := args[0]
			current, err := utils.CurrentBranch()
			if err != nil {
				return err
			}
			if current == name {
				utils.ErrorExit(fmt.Sprintf("Cannot delete current branch %s.", name))
			}
			exists, err := utils.BranchExists(name)
			if err != nil {
				return err
			}
			if !exists {
				utils.ErrorExit(fmt.Sprintf("Branch '%s' does not exist.", name))
			}
			return utils.Git(false, "branch", "-d", name)
		},
	}
}
