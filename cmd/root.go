package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func findGitRepos() ([]string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	gitDir := filepath.Join(homeDir, "Git")
	var repos []string

	err = filepath.Walk(gitDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		relPath, _ := filepath.Rel(gitDir, path)
		if strings.Count(relPath, string(os.PathSeparator)) > 3 {
			return filepath.SkipDir
		}

		if info.IsDir() && info.Name() == ".git" {
			repos = append(repos, filepath.Dir(path))
			return filepath.SkipDir
		}

		return nil
	})

	return repos, err
}

func getRepoNames() []string {
	repos, err := findGitRepos()
	if err != nil {
		return []string{}
	}

	var names []string
	for _, repo := range repos {
		names = append(names, filepath.Base(repo))
	}
	return names
}

var RootCmd = &cobra.Command{
	Use:   "repo-switcher [repo-name]",
	Short: "Switch to a git repository",
	Args:  cobra.ExactArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getRepoNames(), cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		repos, err := findGitRepos()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		targetName := args[0]
		for _, repo := range repos {
			if filepath.Base(repo) == targetName {
				fmt.Printf("cd %s\n", repo)
				return
			}
		}

		fmt.Printf("Repository '%s' not found\n", targetName)
		os.Exit(1)
	},
}
