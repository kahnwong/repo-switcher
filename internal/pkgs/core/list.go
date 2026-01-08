package core

import (
	"os"
	"path/filepath"
	"strings"
)

func listGitRepos() ([]string, error) {
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
