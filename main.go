package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/trypsynth/qit/commands"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "qit",
		Short: "Qit - Quin's tiny Git helper.",
		Long:  "Qit - Quin's tiny Git helper.\nUsage: qit <command> [<args>...]",
	}
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.SetHelpTemplate(`Qit - Quin's tiny Git helper.
Usage: qit <command> [<args>...]
Available commands:
{{range .Commands}}{{if (not .IsAvailableCommand)}}{{else}}  {{.Use}}: {{.Short}}.
{{end}}{{end}}`)
	rootCmd.AddCommand(commands.NewAcpCommand())
	rootCmd.AddCommand(commands.NewAmendCommand())
	rootCmd.AddCommand(commands.NewCpCommand())
	rootCmd.AddCommand(commands.NewDbCommand())
	rootCmd.AddCommand(commands.NewIgnoreCommand())
	rootCmd.AddCommand(commands.NewLastCommand())
	rootCmd.AddCommand(commands.NewLicenseCommand())
	rootCmd.AddCommand(commands.NewLogCommand())
	rootCmd.AddCommand(commands.NewNbCommand())
	rootCmd.AddCommand(commands.NewNewCommand())
	rootCmd.AddCommand(commands.NewResetCommand())
	rootCmd.AddCommand(commands.NewStatusCommand())
	rootCmd.AddCommand(commands.NewUndoCommand())
	rootCmd.SilenceUsage = true
	if err := rootCmd.Execute(); err != nil {
		if cmd, _, _ := rootCmd.Find(os.Args[1:]); cmd == rootCmd && len(os.Args) > 1 {
			fmt.Fprintf(os.Stderr, "Unknown command: %s.\n", os.Args[1])
			fmt.Fprintln(os.Stderr, "Use 'qit help' to see available commands.")
		} else {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(1)
	}
}
