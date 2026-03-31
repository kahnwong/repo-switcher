package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	cliBase "github.com/kahnwong/cli-base"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Paths []string `yaml:"paths"`
}

var AppConfigBasePath string
var AppConfig *Config

var ReposMap map[string]string
var ReposName []string

func init() {
	// Set log level
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Initialize config path
	var err error
	AppConfigBasePath, err = cliBase.ExpandHome("~/.config/repo-switcher")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to expand config path")
	}

	// Initialize cache file path
	cacheFilePath = filepath.Join(AppConfigBasePath, cacheFileName)

	// Read config file
	AppConfig, err = cliBase.ReadYaml[Config](fmt.Sprintf("%s/config.yaml", AppConfigBasePath))
	if err != nil {
		if isTestMode() {
			log.Warn().Err(err).Msg("failed to read config file")
			return
		}
		log.Fatal().Err(err).Msg("failed to read config file")
	}

	if AppConfig == nil {
		if isTestMode() {
			log.Warn().Msg("skipping repo initialization: config not loaded")
			return
		}
		log.Fatal().Msg("config not loaded")
	}

	repos, err := listGitReposWithCache(AppConfig.Paths, false)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to list git repos")
	}

	ReposMap = createGitFolderMap(repos)
	ReposName = getReposName(ReposMap)
}

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

func isTestMode() bool {
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "-test.") || strings.HasSuffix(arg, ".test") {
			return true
		}
	}
	return false
}

// entrypoint - for force refresh
func RefreshCache() error {
	repos, err := listGitReposWithCache(AppConfig.Paths, true)
	if err != nil {
		return err
	}

	ReposMap = createGitFolderMap(repos)
	ReposName = getReposName(ReposMap)
	return nil
}
