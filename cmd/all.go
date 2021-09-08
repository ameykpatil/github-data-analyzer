package cmd

import (
	"github.com/ameykpatil/github-data-analyzer/db"
	"github.com/ameykpatil/github-data-analyzer/domain/repo"
	"github.com/ameykpatil/github-data-analyzer/domain/user"
	"github.com/ameykpatil/github-data-analyzer/service"
	"github.com/spf13/cobra"
)

// NewAllCmd command to get top users & repos
func NewAllCmd() *cobra.Command {
	allCmd := &cobra.Command{
		Use:   "all",
		Short: "Get top users & repos",
		RunE:  getTopUsersAndRepos,
	}

	allCmd.Flags().StringP("path", "p", "", "path of the directory where the data files are")
	allCmd.Flags().Uint32P("limit", "l", 10, "number of users to return")

	return allCmd
}

func getTopUsersAndRepos(cmd *cobra.Command, args []string) error {
	// get & verify flags
	path, err := cmd.Flags().GetString("path")
	if err != nil {
		return err
	}
	limit, err := cmd.Flags().GetUint32("limit")
	if err != nil {
		return err
	}

	// initialise dependencies
	dataStore, err := db.NewDataStore(path)
	if err != nil {
		return err
	}
	eventHandler := service.NewEventHandler(dataStore)
	userAnalyzer := user.NewAnalyzer(*eventHandler)
	repoAnalyzer := repo.NewAnalyzer(*eventHandler)

	// get top users by passing custom sort function
	users := userAnalyzer.GetTopUsers(limit, func(ui, uj user.User) bool {
		if ui.EventTypeCount["PullRequestEvent"] == uj.EventTypeCount["PullRequestEvent"] {
			return ui.CommitCount > uj.CommitCount
		} else if ui.EventTypeCount["PullRequestEvent"] > uj.EventTypeCount["PullRequestEvent"] {
			return true
		}
		return false
	})

	// get top repos by commits by passing custom sort function
	reposByCommits := repoAnalyzer.GetTopRepos(limit, func(ri, rj repo.Repo) bool {
		return ri.CommitCount > rj.CommitCount
	})

	// get top repos by watch events by passing custom sort function
	reposByWatchEvents := repoAnalyzer.GetTopRepos(limit, func(ri, rj repo.Repo) bool {
		return ri.EventTypeCount["WatchEvent"] > rj.EventTypeCount["WatchEvent"]
	})

	// print the results in readable format
	printUsers(users, limit, []string{"PullRequestEvent", "Commits"})
	printRepos(reposByCommits, limit, []string{"Commits"})
	printRepos(reposByWatchEvents, limit, []string{"WatchEvent"})

	return nil
}
