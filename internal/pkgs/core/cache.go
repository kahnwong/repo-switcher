package core

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

type RepoCache struct {
	Repos     []string  `json:"repos"`
	Timestamp time.Time `json:"timestamp"`
	PathsHash string    `json:"paths_hash"`
}

const (
	cacheTTL      = 24 * time.Hour
	cacheFileName = "repos-cache.json"
)

var cacheFilePath = filepath.Join(AppConfigBasePath, cacheFileName)

// hashPaths creates a hash of the paths to detect config changes
func hashPaths(paths []string) string {
	h := sha256.New()
	h.Write([]byte(strings.Join(paths, "|")))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// readCache reads the cache from disk
func readCache() (*RepoCache, error) {
	data, err := os.ReadFile(cacheFilePath)
	if err != nil {
		return nil, err
	}

	var cache RepoCache
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, err
	}

	return &cache, nil
}

// writeCache writes the cache to disk
func writeCache(repos []string, paths []string) error {
	cache := RepoCache{
		Repos:     repos,
		Timestamp: time.Now(),
		PathsHash: hashPaths(paths),
	}

	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return err
	}

	// Ensure the directory exists
	if err := os.MkdirAll(AppConfigBasePath, 0755); err != nil {
		return err
	}

	return os.WriteFile(cacheFilePath, data, 0644)
}

// isCacheValid checks if the cache is still valid
func isCacheValid(cache *RepoCache, paths []string) bool {
	// Check if cache is too old
	if time.Since(cache.Timestamp) > cacheTTL {
		log.Debug().Msg("cache expired")
		return false
	}

	// Check if paths have changed
	currentHash := hashPaths(paths)
	if cache.PathsHash != currentHash {
		log.Debug().Msg("paths configuration changed")
		return false
	}

	return true
}

// listGitReposWithCache returns git repos using cache when possible
func listGitReposWithCache(paths []string, forceRefresh bool) ([]string, error) {
	// If force refresh is requested, skip cache
	if !forceRefresh {
		cache, err := readCache()
		if err == nil && isCacheValid(cache, paths) {
			log.Debug().Msg("using cached repository list")
			return cache.Repos, nil
		}
		if err != nil {
			log.Debug().Err(err).Msg("failed to read cache")
		}
	}

	// Cache miss or invalid - scan directories
	log.Debug().Msg("scanning directories for git repositories")
	repos, err := listGitRepos(paths)
	if err != nil {
		return nil, err
	}

	// Write to cache
	if err := writeCache(repos, paths); err != nil {
		log.Warn().Err(err).Msg("failed to write cache")
		// Don't fail if cache write fails, just continue
	}

	return repos, nil
}
