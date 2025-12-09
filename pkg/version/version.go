package version

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Version is set by goreleaser or build script via ldflags
var (
	Version = "v0.0.2"
	Commit  = "none"
	Date    = "unknown"
)

// GitHubTag represents a tag from GitHub API
type GitHubTag struct {
	Name string `json:"name"`
}

// CheckForUpdate checks GitHub for a newer version and prints a message if available.
// Errors are handled silently - returns without printing if check fails.
func CheckForUpdate() {
	latest, err := getLatestVersion()
	if err != nil {
		return
	}

	if latest != "" && latest != Version && isNewer(latest, Version) {
		fmt.Printf("\nA new version of coolpack is available: %s (current: %s)\n", latest, Version)
		fmt.Println("Download: https://github.com/coollabsio/coolpack/releases/latest")
	}
}

// getLatestVersion fetches the latest release tag from GitHub
func getLatestVersion() (string, error) {
	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get("https://api.github.com/repos/coollabsio/coolpack/tags?per_page=10")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("github api returned %d", resp.StatusCode)
	}

	var tags []GitHubTag
	if err := json.NewDecoder(resp.Body).Decode(&tags); err != nil {
		return "", err
	}

	// Find the latest version tag (starts with 'v')
	for _, tag := range tags {
		if strings.HasPrefix(tag.Name, "v") {
			return tag.Name, nil
		}
	}

	return "", nil
}

// isNewer compares two semver strings and returns true if latest > current
func isNewer(latest, current string) bool {
	// Strip 'v' prefix
	latest = strings.TrimPrefix(latest, "v")
	current = strings.TrimPrefix(current, "v")

	latestParts := strings.Split(latest, ".")
	currentParts := strings.Split(current, ".")

	// Compare major.minor.patch
	for i := 0; i < 3; i++ {
		var latestNum, currentNum int
		if i < len(latestParts) {
			fmt.Sscanf(latestParts[i], "%d", &latestNum)
		}
		if i < len(currentParts) {
			fmt.Sscanf(currentParts[i], "%d", &currentNum)
		}

		if latestNum > currentNum {
			return true
		}
		if latestNum < currentNum {
			return false
		}
	}

	return false
}
