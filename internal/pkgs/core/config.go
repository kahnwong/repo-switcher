package core

import "path/filepath"

var ReposMap map[string]string
var ReposName []string

func createGitFolderMap() map[string]string {
	repos, err := listGitRepos()
	if err != nil {
		return map[string]string{}
	}

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
	ReposMap = createGitFolderMap()
	ReposName = getReposName(ReposMap)
}
