package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/trypsynth/qit/utils"
)

func NewSummaryCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "summary",
		Short: "show commit summary by author",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			output, err := utils.GitOutput("log", "--pretty=%an")
			if err != nil {
				return err
			}
			authorCounts := make(map[string]int)
			lines := strings.Split(output, "\n")
			for _, author := range lines {
				author = strings.TrimSpace(author)
				if author != "" {
					authorCounts[author]++
				}
			}
			type authorCount struct {
				name  string
				count int
			}
			var sorted []authorCount
			for name, count := range authorCounts {
				sorted = append(sorted, authorCount{name, count})
			}
			for i := 0; i < len(sorted); i++ {
				for j := i + 1; j < len(sorted); j++ {
					if sorted[j].count > sorted[i].count {
						sorted[i], sorted[j] = sorted[j], sorted[i]
					}
				}
			}
			for _, ac := range sorted {
				commitWord := "commits"
				if ac.count == 1 {
					commitWord = "commit"
				}
				fmt.Printf("%s has made %d %s.\n", ac.name, ac.count, commitWord)
			}
			return nil
		},
	}
}
