package cmd

import (
	"fmt"
	"os"

	"github.com/kahnwong/repo-switcher/internal/pkgs/core"
	"github.com/spf13/cobra"
)

var reposMap = core.ReposMap
var reposName = core.ReposName

var RootCmd = &cobra.Command{
	Use:          "repo-switcher [repo-name]",
	Short:        "Switch to a git repository",
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return reposName, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		repoName := args[0]

		if fullPath, exists := reposMap[repoName]; exists {
			fmt.Println(fullPath)
			os.Exit(0)
		}

		fmt.Printf("Repository '%s' not found\n", repoName)
		os.Exit(1)
	},
}
