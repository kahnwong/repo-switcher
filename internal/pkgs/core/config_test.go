package core

import (
	"reflect"
	"testing"
)

func TestCreateGitFolderMap(t *testing.T) {
	tests := []struct {
		name     string
		repos    []string
		expected map[string]string
	}{
		{
			name:     "empty repos",
			repos:    []string{},
			expected: map[string]string{},
		},
		{
			name: "single repo",
			repos: []string{
				"/home/user/projects/repo1",
			},
			expected: map[string]string{
				"repo1": "/home/user/projects/repo1",
			},
		},
		{
			name: "multiple repos",
			repos: []string{
				"/home/user/projects/repo1",
				"/home/user/projects/repo2",
				"/var/www/myapp",
			},
			expected: map[string]string{
				"repo1": "/home/user/projects/repo1",
				"repo2": "/home/user/projects/repo2",
				"myapp": "/var/www/myapp",
			},
		},
		{
			name: "repos with same basename",
			repos: []string{
				"/home/user/projects/repo1",
				"/home/user/work/repo1",
			},
			expected: map[string]string{
				"repo1": "/home/user/work/repo1", // last one wins
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := createGitFolderMap(tt.repos)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("createGitFolderMap() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetReposName(t *testing.T) {
	tests := []struct {
		name     string
		reposMap map[string]string
		expected int // we check length since map iteration order is random
	}{
		{
			name:     "empty map",
			reposMap: map[string]string{},
			expected: 0,
		},
		{
			name: "single entry",
			reposMap: map[string]string{
				"repo1": "/home/user/projects/repo1",
			},
			expected: 1,
		},
		{
			name: "multiple entries",
			reposMap: map[string]string{
				"repo1": "/home/user/projects/repo1",
				"repo2": "/home/user/projects/repo2",
				"repo3": "/home/user/projects/repo3",
			},
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getReposName(tt.reposMap)
			if len(result) != tt.expected {
				t.Errorf("getReposName() returned %d items, want %d", len(result), tt.expected)
			}

			// Verify all keys are present
			for key := range tt.reposMap {
				found := false
				for _, name := range result {
					if name == key {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("getReposName() missing key %s", key)
				}
			}
		})
	}
}
