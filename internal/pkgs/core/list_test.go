package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestListGitRepos(t *testing.T) {
	// Create a temporary directory structure for testing
	tempDir := t.TempDir()

	// Create test directory structure with git repos
	testRepos := []struct {
		path      string
		isGitRepo bool
	}{
		{"repo1/.git", true},
		{"repo2/.git", true},
		{"nested/repo3/.git", true},
		{"nested/deep/repo4/.git", true},
		{"not-a-repo", false},
		{"nested/not-a-repo", false},
		// This should be skipped due to depth limit (> 3 levels)
		{"level1/level2/level3/level4/.git", true},
	}

	for _, repo := range testRepos {
		fullPath := filepath.Join(tempDir, repo.path)
		err := os.MkdirAll(fullPath, 0755)
		if err != nil {
			t.Fatalf("failed to create test directory %s: %v", fullPath, err)
		}
	}

	// Test listing git repos
	repos, err := listGitRepos([]string{tempDir})
	if err != nil {
		t.Fatalf("listGitRepos() error = %v", err)
	}

	// Expected repos (excluding the one that's too deep)
	expectedRepos := []string{
		filepath.Join(tempDir, "repo1"),
		filepath.Join(tempDir, "repo2"),
		filepath.Join(tempDir, "nested/repo3"),
		filepath.Join(tempDir, "nested/deep/repo4"),
	}

	if len(repos) != len(expectedRepos) {
		t.Errorf("listGitRepos() found %d repos, want %d", len(repos), len(expectedRepos))
		t.Logf("Found repos: %v", repos)
		t.Logf("Expected repos: %v", expectedRepos)
	}

	// Verify each expected repo is in the results
	for _, expected := range expectedRepos {
		found := false
		for _, repo := range repos {
			if repo == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("listGitRepos() missing expected repo: %s", expected)
		}
	}
}

func TestListGitReposEmptyDirectory(t *testing.T) {
	tempDir := t.TempDir()

	repos, err := listGitRepos([]string{tempDir})
	if err != nil {
		t.Fatalf("listGitRepos() error = %v", err)
	}

	if len(repos) != 0 {
		t.Errorf("listGitRepos() found %d repos in empty directory, want 0", len(repos))
	}
}

func TestListGitReposMultiplePaths(t *testing.T) {
	// Create two separate temp directories
	tempDir1 := t.TempDir()
	tempDir2 := t.TempDir()

	// Create git repos in first directory
	err := os.MkdirAll(filepath.Join(tempDir1, "repo1/.git"), 0755)
	if err != nil {
		t.Fatalf("failed to create test directory: %v", err)
	}

	// Create git repos in second directory
	err = os.MkdirAll(filepath.Join(tempDir2, "repo2/.git"), 0755)
	if err != nil {
		t.Fatalf("failed to create test directory: %v", err)
	}

	// Test listing git repos from both paths
	repos, err := listGitRepos([]string{tempDir1, tempDir2})
	if err != nil {
		t.Fatalf("listGitRepos() error = %v", err)
	}

	if len(repos) != 2 {
		t.Errorf("listGitRepos() found %d repos, want 2", len(repos))
	}

	// Verify both repos are found
	expectedRepos := []string{
		filepath.Join(tempDir1, "repo1"),
		filepath.Join(tempDir2, "repo2"),
	}

	for _, expected := range expectedRepos {
		found := false
		for _, repo := range repos {
			if repo == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("listGitRepos() missing expected repo: %s", expected)
		}
	}
}

func TestListGitReposNonExistentPath(t *testing.T) {
	// Test with a path that doesn't exist
	nonExistentPath := "/this/path/does/not/exist/hopefully"

	repos, err := listGitRepos([]string{nonExistentPath})

	// The function should handle the error gracefully and return empty list
	if err != nil {
		t.Logf("listGitRepos() returned error for non-existent path: %v", err)
	}

	// Should return empty list or handle gracefully
	if repos == nil {
		repos = []string{}
	}

	t.Logf("listGitRepos() found %d repos for non-existent path", len(repos))
}

func TestListGitReposDepthLimit(t *testing.T) {
	tempDir := t.TempDir()

	// Create repos at different depths
	// Note: depth is counted by path separators in the relative path
	testCases := []struct {
		path          string
		shouldBeFound bool
	}{
		{"repo1/.git", true},                                     // 0 separators
		{"level1/repo2/.git", true},                              // 1 separator
		{"level1/level2/repo3/.git", true},                       // 2 separators
		{"level1/level2/level3/repo4/.git", false},               // 3 separators - at the limit, .git adds one more
		{"level1/level2/level3/level4/repo5/.git", false},        // 4 separators - should be skipped
		{"level1/level2/level3/level4/level5/repo6/.git", false}, // 5 separators - should be skipped
	}

	for _, tc := range testCases {
		fullPath := filepath.Join(tempDir, tc.path)
		err := os.MkdirAll(fullPath, 0755)
		if err != nil {
			t.Fatalf("failed to create test directory %s: %v", fullPath, err)
		}
	}

	repos, err := listGitRepos([]string{tempDir})
	if err != nil {
		t.Fatalf("listGitRepos() error = %v", err)
	}

	// Count how many should be found
	expectedCount := 0
	for _, tc := range testCases {
		if tc.shouldBeFound {
			expectedCount++
		}
	}

	if len(repos) != expectedCount {
		t.Errorf("listGitRepos() found %d repos, want %d", len(repos), expectedCount)
		t.Logf("Found repos: %v", repos)
	}

	// Verify repos that should be found are present
	for _, tc := range testCases {
		if tc.shouldBeFound {
			expectedPath := filepath.Join(tempDir, filepath.Dir(tc.path))
			found := false
			for _, repo := range repos {
				if repo == expectedPath {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("listGitRepos() missing expected repo at depth: %s", tc.path)
			}
		}
	}

	// Verify repos that should NOT be found are absent
	for _, tc := range testCases {
		if !tc.shouldBeFound {
			unexpectedPath := filepath.Join(tempDir, filepath.Dir(tc.path))
			for _, repo := range repos {
				if repo == unexpectedPath {
					t.Errorf("listGitRepos() found repo that should be skipped due to depth: %s", tc.path)
				}
			}
		}
	}
}

func TestListGitReposNestedGitRepos(t *testing.T) {
	tempDir := t.TempDir()

	// Create nested git repos (repo within a repo)
	err := os.MkdirAll(filepath.Join(tempDir, "parent-repo/.git"), 0755)
	if err != nil {
		t.Fatalf("failed to create test directory: %v", err)
	}

	err = os.MkdirAll(filepath.Join(tempDir, "parent-repo/nested-repo/.git"), 0755)
	if err != nil {
		t.Fatalf("failed to create test directory: %v", err)
	}

	repos, err := listGitRepos([]string{tempDir})
	if err != nil {
		t.Fatalf("listGitRepos() error = %v", err)
	}

	// Should find both repos
	expectedRepos := []string{
		filepath.Join(tempDir, "parent-repo"),
		filepath.Join(tempDir, "parent-repo/nested-repo"),
	}

	if len(repos) != len(expectedRepos) {
		t.Errorf("listGitRepos() found %d repos, want %d", len(repos), len(expectedRepos))
	}

	for _, expected := range expectedRepos {
		found := false
		for _, repo := range repos {
			if repo == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("listGitRepos() missing expected repo: %s", expected)
		}
	}
}
