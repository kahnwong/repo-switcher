package core

import (
	"os"
	"path/filepath"
	"strings"

	cli_base "github.com/kahnwong/cli-base"
)

func listGitRepos(paths []string) ([]string, error) {
	var repos []string
	var err error

	for _, path := range paths {
		gitDir, _ := cli_base.ExpandHome(path)

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
	}

	return repos, err
}
