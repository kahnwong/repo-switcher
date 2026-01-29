package core

import (
	"os"
	"path/filepath"

	"fmt"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	cliBase "github.com/kahnwong/cli-base"
)

func init() {
	// Set log level to info before any logging occurs
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

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
	repos, err := listGitReposWithCache(AppConfig.Paths, false)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to list git repos")
	}

	ReposMap = createGitFolderMap(repos)
	ReposName = getReposName(ReposMap)
}

// RefreshCache forces a refresh of the repository cache
func RefreshCache() error {
	repos, err := listGitReposWithCache(AppConfig.Paths, true)
	if err != nil {
		return err
	}

	ReposMap = createGitFolderMap(repos)
	ReposName = getReposName(ReposMap)
	return nil
}
