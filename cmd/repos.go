package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ameykpatil/github-data-analyzer/db"
	"github.com/ameykpatil/github-data-analyzer/domain/repo"
	"github.com/ameykpatil/github-data-analyzer/service"
	"github.com/spf13/cobra"
)

// NewReposCmd command to get top repos
func NewReposCmd() *cobra.Command {
	reposCmd := &cobra.Command{
		Use:   "repos",
		Short: "Get the top repos",
		RunE:  getTopRepos,
	}

	reposCmd.Flags().StringP("path", "p", "", "path of the directory where the data files are")
	reposCmd.Flags().Uint32P("limit", "l", 10, "number of users to return")
	reposCmd.Flags().StringP("sort", "s", "commits", "field to sort by")

	return reposCmd
}

func getTopRepos(cmd *cobra.Command, args []string) error {
	// get & verify the flags
	path, err := cmd.Flags().GetString("path")
	if err != nil {
		return err
	}
	limit, err := cmd.Flags().GetUint32("limit")
	if err != nil {
		return err
	}
	sortField, err := cmd.Flags().GetString("sort")
	if err != nil {
		return err
	}

	// create sort function based on sort field
	var fn func(ri, rj repo.Repo) bool
	switch sortField {
	case "Commits":
		fn = func(ri, rj repo.Repo) bool {
			return ri.CommitCount > rj.CommitCount
		}
	default:
		if strings.Contains(sortField, "Event") {
			fn = func(ri, rj repo.Repo) bool {
				return ri.EventTypeCount[sortField] > rj.EventTypeCount[sortField]
			}
		} else {
			return errors.New("invalid sort field " + sortField)
		}
	}

	// initialise dependencies
	dataStore, err := db.NewDataStore(path)
	if err != nil {
		return err
	}
	eventHandler := service.NewEventHandler(dataStore)
	repoAnalyzer := repo.NewAnalyzer(*eventHandler)

	// get the top repos
	repos := repoAnalyzer.GetTopRepos(limit, fn)

	// print the result in readable format
	printRepos(repos, limit, []string{sortField})

	return nil
}

// printRepos print repos in readable format
func printRepos(repos []repo.Repo, limit uint32, sortFields []string) {
	var str strings.Builder
	for _, repo := range repos {
		for _, sortField := range sortFields {
			if sortField == "Commits" {
				fmt.Fprintf(&str, "%s:%d ", sortField, repo.CommitCount)
			} else if strings.Contains(sortField, "Event") {
				fmt.Fprintf(&str, "%s:%d ", sortField, repo.EventTypeCount[sortField])
			}
		}
		fmt.Fprintf(&str, "ID:%s Name:%s \n", repo.ID, repo.Name)
	}
	fmt.Printf("Top %d Repos by %v \n --- \n%s --- \n", limit, sortFields, str.String())
}
