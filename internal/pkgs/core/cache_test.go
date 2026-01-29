package core

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestHashPaths(t *testing.T) {
	tests := []struct {
		name     string
		paths    []string
		expected string
	}{
		{
			name:     "empty paths",
			paths:    []string{},
			expected: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
		{
			name:     "single path",
			paths:    []string{"/home/user/projects"},
			expected: hashPaths([]string{"/home/user/projects"}), // compute expected
		},
		{
			name:  "multiple paths",
			paths: []string{"/home/user/projects", "/var/www"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hashPaths(tt.paths)
			if len(result) != 64 { // SHA256 produces 64 hex characters
				t.Errorf("hashPaths() returned hash of length %d, want 64", len(result))
			}

			// Test consistency - same input should produce same hash
			result2 := hashPaths(tt.paths)
			if result != result2 {
				t.Errorf("hashPaths() not consistent: %s != %s", result, result2)
			}
		})
	}

	// Test that different paths produce different hashes
	t.Run("different paths produce different hashes", func(t *testing.T) {
		hash1 := hashPaths([]string{"/path1"})
		hash2 := hashPaths([]string{"/path2"})
		if hash1 == hash2 {
			t.Error("hashPaths() produced same hash for different paths")
		}
	})

	// Test that order matters
	t.Run("order matters", func(t *testing.T) {
		hash1 := hashPaths([]string{"/path1", "/path2"})
		hash2 := hashPaths([]string{"/path2", "/path1"})
		if hash1 == hash2 {
			t.Error("hashPaths() produced same hash for different order")
		}
	})
}

func TestWriteAndReadCache(t *testing.T) {
	// Save original cache file path and restore after test
	originalCachePath := cacheFilePath
	defer func() { cacheFilePath = originalCachePath }()

	// Create temp directory for test
	tempDir := t.TempDir()
	cacheFilePath = filepath.Join(tempDir, cacheFileName)

	repos := []string{
		"/home/user/projects/repo1",
		"/home/user/projects/repo2",
	}
	paths := []string{"/home/user/projects"}

	// Test writing cache
	err := writeCache(repos, paths)
	if err != nil {
		t.Fatalf("writeCache() error = %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(cacheFilePath); os.IsNotExist(err) {
		t.Fatal("cache file was not created")
	}

	// Test reading cache
	cache, err := readCache()
	if err != nil {
		t.Fatalf("readCache() error = %v", err)
	}

	// Verify cache contents
	if len(cache.Repos) != len(repos) {
		t.Errorf("readCache() returned %d repos, want %d", len(cache.Repos), len(repos))
	}

	for i, repo := range repos {
		if cache.Repos[i] != repo {
			t.Errorf("readCache() repo[%d] = %s, want %s", i, cache.Repos[i], repo)
		}
	}

	if cache.PathsHash != hashPaths(paths) {
		t.Errorf("readCache() PathsHash = %s, want %s", cache.PathsHash, hashPaths(paths))
	}

	// Verify timestamp is recent
	if time.Since(cache.Timestamp) > time.Minute {
		t.Error("readCache() timestamp is not recent")
	}
}

func TestReadCacheNonExistent(t *testing.T) {
	// Save original cache file path and restore after test
	originalCachePath := cacheFilePath
	defer func() { cacheFilePath = originalCachePath }()

	// Point to non-existent file
	cacheFilePath = filepath.Join(t.TempDir(), "nonexistent.json")

	_, err := readCache()
	if err == nil {
		t.Error("readCache() expected error for non-existent file, got nil")
	}
}

func TestReadCacheInvalidJSON(t *testing.T) {
	// Save original cache file path and restore after test
	originalCachePath := cacheFilePath
	defer func() { cacheFilePath = originalCachePath }()

	tempDir := t.TempDir()
	cacheFilePath = filepath.Join(tempDir, cacheFileName)

	// Write invalid JSON
	err := os.WriteFile(cacheFilePath, []byte("invalid json"), 0644)
	if err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	_, err = readCache()
	if err == nil {
		t.Error("readCache() expected error for invalid JSON, got nil")
	}
}

func TestIsCacheValid(t *testing.T) {
	paths := []string{"/home/user/projects"}
	pathsHash := hashPaths(paths)

	tests := []struct {
		name     string
		cache    *RepoCache
		paths    []string
		expected bool
	}{
		{
			name: "valid cache",
			cache: &RepoCache{
				Repos:     []string{"/home/user/projects/repo1"},
				Timestamp: time.Now(),
				PathsHash: pathsHash,
			},
			paths:    paths,
			expected: true,
		},
		{
			name: "expired cache",
			cache: &RepoCache{
				Repos:     []string{"/home/user/projects/repo1"},
				Timestamp: time.Now().Add(-25 * time.Hour), // older than cacheTTL
				PathsHash: pathsHash,
			},
			paths:    paths,
			expected: false,
		},
		{
			name: "paths changed",
			cache: &RepoCache{
				Repos:     []string{"/home/user/projects/repo1"},
				Timestamp: time.Now(),
				PathsHash: hashPaths([]string{"/different/path"}),
			},
			paths:    paths,
			expected: false,
		},
		{
			name: "both expired and paths changed",
			cache: &RepoCache{
				Repos:     []string{"/home/user/projects/repo1"},
				Timestamp: time.Now().Add(-25 * time.Hour),
				PathsHash: hashPaths([]string{"/different/path"}),
			},
			paths:    paths,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isCacheValid(tt.cache, tt.paths)
			if result != tt.expected {
				t.Errorf("isCacheValid() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestWriteCacheCreatesDirectory(t *testing.T) {
	// Save original values and restore after test
	originalCachePath := cacheFilePath
	originalAppConfigBasePath := AppConfigBasePath
	defer func() {
		cacheFilePath = originalCachePath
		AppConfigBasePath = originalAppConfigBasePath
	}()

	// Create temp directory
	tempDir := t.TempDir()
	nestedDir := filepath.Join(tempDir, "nested", "config")
	AppConfigBasePath = nestedDir
	cacheFilePath = filepath.Join(nestedDir, cacheFileName)

	repos := []string{"/home/user/projects/repo1"}
	paths := []string{"/home/user/projects"}

	// Directory should not exist yet
	if _, err := os.Stat(nestedDir); !os.IsNotExist(err) {
		t.Fatal("nested directory should not exist yet")
	}

	// Write cache should create the directory
	err := writeCache(repos, paths)
	if err != nil {
		t.Fatalf("writeCache() error = %v", err)
	}

	// Verify directory was created
	if _, err := os.Stat(nestedDir); os.IsNotExist(err) {
		t.Error("writeCache() did not create directory")
	}

	// Verify file was created
	if _, err := os.Stat(cacheFilePath); os.IsNotExist(err) {
		t.Error("writeCache() did not create cache file")
	}
}

func TestCacheJSONFormat(t *testing.T) {
	// Save original cache file path and restore after test
	originalCachePath := cacheFilePath
	defer func() { cacheFilePath = originalCachePath }()

	tempDir := t.TempDir()
	cacheFilePath = filepath.Join(tempDir, cacheFileName)

	repos := []string{"/home/user/projects/repo1"}
	paths := []string{"/home/user/projects"}

	err := writeCache(repos, paths)
	if err != nil {
		t.Fatalf("writeCache() error = %v", err)
	}

	// Read raw JSON and verify structure
	data, err := os.ReadFile(cacheFilePath)
	if err != nil {
		t.Fatalf("failed to read cache file: %v", err)
	}

	var rawCache map[string]interface{}
	err = json.Unmarshal(data, &rawCache)
	if err != nil {
		t.Fatalf("cache file is not valid JSON: %v", err)
	}

	// Verify expected fields exist
	if _, ok := rawCache["repos"]; !ok {
		t.Error("cache JSON missing 'repos' field")
	}
	if _, ok := rawCache["timestamp"]; !ok {
		t.Error("cache JSON missing 'timestamp' field")
	}
	if _, ok := rawCache["paths_hash"]; !ok {
		t.Error("cache JSON missing 'paths_hash' field")
	}
}
