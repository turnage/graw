package graw

import (
	"strings"
)

// subredditsPath returns a path to listing with all requested subreddits.
func subredditsPath(subs []string) string {
	return strings.Join(
		[]string{
			"/r",
			strings.Join(subs, "+"),
			"new",
		}, "/",
	)
}

// userPaths returns paths to the user accounts specified.
func userPaths(users []string) []string {
	paths := make([]string, len(users))
	for i, user := range users {
		paths[i] = strings.Join([]string{"/u", user}, "/")
	}
	return paths
}

// logPathsOut transforms paths into one that explicitly requests the raw json
// at the endpoint, because logged out paths by default provide user-facing
// html, css, etc.
func logPathsOut(paths []string) []string {
	loggedOutPaths := make([]string, len(paths))
	for i, path := range paths {
		loggedOutPaths[i] = path + ".json"
	}
	return loggedOutPaths
}
