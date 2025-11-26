package commands

import (
	"fmt"
	"sort"
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
			sort.Slice(sorted, func(i, j int) bool {
				return sorted[i].count > sorted[j].count
			})
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
