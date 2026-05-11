package core

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	cliBase "github.com/kahnwong/cli-base"
	"github.com/rs/zerolog"
	slogzerolog "github.com/samber/slog-zerolog/v2"
)

type Config struct {
	Paths []string `yaml:"paths"`
}

var AppConfigBasePath string
var AppConfig *Config

var ReposMap map[string]string
var ReposName []string

func init() {
	output := zerolog.ConsoleWriter{Out: os.Stderr}
	logger := zerolog.New(output).Level(zerolog.ErrorLevel).With().Timestamp().Logger()
	slog.SetDefault(slog.New(slogzerolog.Option{Logger: &logger}.NewZerologHandler()))

	// Initialize config path
	var err error
	AppConfigBasePath, err = cliBase.ExpandHome("~/.config/repo-switcher")
	if err != nil {
		slog.Error("failed to expand config path", "error", err)
		os.Exit(1)
	}

	// Initialize cache file path
	cacheFilePath = filepath.Join(AppConfigBasePath, cacheFileName)

	// Read config file
	AppConfig, err = cliBase.ReadYaml[Config](fmt.Sprintf("%s/config.yaml", AppConfigBasePath))
	if err != nil {
		if isTestMode() {
			slog.Warn("failed to read config file", "error", err)
			return
		}
		slog.Error("failed to read config file", "error", err)
		os.Exit(1)
	}

	if AppConfig == nil {
		if isTestMode() {
			slog.Warn("skipping repo initialization: config not loaded")
			return
		}
		slog.Error("config not loaded")
		os.Exit(1)
	}

	repos, err := listGitReposWithCache(AppConfig.Paths, false)
	if err != nil {
		slog.Error("failed to list git repos", "error", err)
		os.Exit(1)
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
