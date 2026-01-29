package cmd

import (
	"fmt"

	"github.com/kahnwong/repo-switcher/internal/pkgs/core"
	"github.com/spf13/cobra"
)

var refreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Refresh the repository cache",
	Long:  "Scans all configured paths and updates the repository cache",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Refreshing repository cache...")

		if err := core.RefreshCache(); err != nil {
			fmt.Printf("Error refreshing cache: %v\n", err)
			return
		}

		fmt.Printf("Cache refreshed successfully. Found %d repositories.\n", len(core.ReposName))
	},
}

func init() {
	RootCmd.AddCommand(refreshCmd)
}
