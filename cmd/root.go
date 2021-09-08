package cmd

import (
	"github.com/spf13/cobra"
)

// NewRootCmd creates the root command
func NewRootCmd(version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "analyzer",
		Version: version,
	}

	cmd.AddCommand(NewAllCmd())
	cmd.AddCommand(NewUsersCmd())
	cmd.AddCommand(NewReposCmd())

	return cmd
}
