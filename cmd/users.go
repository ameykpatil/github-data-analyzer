package cmd

import (
	"fmt"
	"strings"

	"github.com/ameykpatil/github-data-analyzer/db"
	"github.com/ameykpatil/github-data-analyzer/domain/user"
	"github.com/ameykpatil/github-data-analyzer/service"
	"github.com/spf13/cobra"
)

// NewUsersCmd command to get top users
func NewUsersCmd() *cobra.Command {
	usersCmd := &cobra.Command{
		Use:   "users",
		Short: "Get top users",
		RunE:  getTopUsers,
	}

	usersCmd.Flags().StringP("path", "p", "", "path of the directory where the data files are")
	usersCmd.Flags().Uint32P("limit", "l", 10, "number of users to return")
	usersCmd.Flags().StringSliceP("sort", "s", []string{"prs,commits"}, "fields to sort by")

	return usersCmd
}

func getTopUsers(cmd *cobra.Command, args []string) error {
	// get & verify flags
	path, err := cmd.Flags().GetString("path")
	if err != nil {
		return err
	}
	limit, err := cmd.Flags().GetUint32("limit")
	if err != nil {
		return err
	}
	sortFields, err := cmd.Flags().GetStringSlice("sort")
	if err != nil {
		return err
	}

	// create custom sort function
	sortFn, err := getSortFunction(sortFields)
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

	// get the top users
	users := userAnalyzer.GetTopUsers(limit, sortFn)

	// print the result in readable format
	printUsers(users, limit, sortFields)

	return nil

}

// getSortFunction creates multilevel wrapped function based on sortFields
func getSortFunction(sortFields []string) (func(ui, uj user.User) bool, error) {
	sortFn := func(ui, uj user.User) bool {
		return true
	}

	for i := len(sortFields) - 1; i >= 0; i-- {
		field := sortFields[i]
		sortFn = wrap(sortFn, field)
	}

	return sortFn, nil
}

// wrap function wraps the given function with sorting of the given field to return a new function
func wrap(f func(ui, uj user.User) bool, field string) func(ui, uj user.User) bool {
	switch field {
	case "Commits":
		return func(ui, uj user.User) bool {
			if ui.CommitCount == uj.CommitCount {
				return f(ui, uj)
			} else if ui.CommitCount > uj.CommitCount {
				return true
			}
			return false
		}
	default:
		if strings.Contains(field, "Event") {
			return func(ri, rj user.User) bool {
				if ri.EventTypeCount[field] == rj.EventTypeCount[field] {
					return f(ri, rj)
				} else if ri.EventTypeCount[field] > rj.EventTypeCount[field] {
					return true
				}
				return false
			}
		}
		return f
	}
}

// printUsers print users in readable format
func printUsers(users []user.User, limit uint32, sortFields []string) {
	var str strings.Builder
	for _, user := range users {
		for _, sortField := range sortFields {
			if sortField == "Commits" {
				fmt.Fprintf(&str, "%s:%d ", sortField, user.CommitCount)
			} else if strings.Contains(sortField, "Event") {
				fmt.Fprintf(&str, "%s:%d ", sortField, user.EventTypeCount[sortField])
			}
		}
		fmt.Fprintf(&str, "ID:%s Username:%s \n", user.ID, user.Username)
	}

	fmt.Printf("Top %d Users by %v \n --- \n%s --- \n", limit, sortFields, str.String())
}
