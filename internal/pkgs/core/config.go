package core

import (
	"path/filepath"

	"fmt"

	"github.com/rs/zerolog/log"

	cliBase "github.com/kahnwong/cli-base"
)

type Config struct {
	Paths []string `yaml:"paths"`
}

var AppConfigBasePath = cliBase.ExpandHome("~/.config/repo-switcher")
var AppConfig = cliBase.ReadYaml[Config](fmt.Sprintf("%s/config.yaml", AppConfigBasePath)) // init

var ReposMap map[string]string
var ReposName []string

func createGitFolderMap(repos []string) map[string]string {
	folderMap := make(map[string]string)
	for _, repo := range repos {
		folderName := filepath.Base(repo)
		folderMap[folderName] = repo
	}
	return folderMap
}

func getReposName(reposMap map[string]string) []string {
	keys := make([]string, 0, len(reposMap))
	for key := range reposMap {
		keys = append(keys, key)
	}
	return keys
}

func init() {
	repos, err := listGitRepos(AppConfig.Paths)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to list git repos")
	}

	ReposMap = createGitFolderMap(repos)
	ReposName = getReposName(ReposMap)
}
