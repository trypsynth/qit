package commands

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/trypsynth/qit/utils"
)

func NewSummaryCommand() *cobra.Command {
	var showLines bool
	cmd := &cobra.Command{
		Use:   "summary",
		Short: "show commit summary by author",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			type authorStats struct {
				commits   int
				additions int
				deletions int
			}
			stats := make(map[string]*authorStats)
			if showLines {
				output, err := utils.GitOutput("log", "--numstat", "--pretty=tformat:%an")
				if err != nil {
					return err
				}
				var currentAuthor string
				for _, line := range strings.Split(output, "\n") {
					line = strings.TrimSpace(line)
					if line == "" {
						continue
					}
					parts := strings.SplitN(line, "\t", 3)
					if len(parts) == 3 {
						if stats[currentAuthor] == nil {
							continue
						}
						add, err1 := strconv.Atoi(parts[0])
						del, err2 := strconv.Atoi(parts[1])
						if err1 != nil || err2 != nil {
							continue // binary file
						}
						stats[currentAuthor].additions += add
						stats[currentAuthor].deletions += del
					} else {
						currentAuthor = line
						if stats[currentAuthor] == nil {
							stats[currentAuthor] = &authorStats{}
						}
						stats[currentAuthor].commits++
					}
				}
			} else {
				output, err := utils.GitOutput("log", "--pretty=%an")
				if err != nil {
					return err
				}
				for _, author := range strings.Split(output, "\n") {
					author = strings.TrimSpace(author)
					if author == "" {
						continue
					}
					if stats[author] == nil {
						stats[author] = &authorStats{}
					}
					stats[author].commits++
				}
			}
			type authorEntry struct {
				name string
				*authorStats
			}
			var sorted []authorEntry
			for name, s := range stats {
				sorted = append(sorted, authorEntry{name, s})
			}
			sort.Slice(sorted, func(i, j int) bool {
				return sorted[i].commits > sorted[j].commits
			})
			for _, e := range sorted {
				commitWord := "commits"
				if e.commits == 1 {
					commitWord = "commit"
				}
				if showLines {
					fmt.Printf("%s has made %d %s (+%d/-%d lines).\n", e.name, e.commits, commitWord, e.additions, e.deletions)
				} else {
					fmt.Printf("%s has made %d %s.\n", e.name, e.commits, commitWord)
				}
			}
			return nil
		},
	}
	cmd.Flags().BoolVarP(&showLines, "lines", "l", false, "include lines added/removed per author")
	return cmd
}
