package git

import (
	"fmt"
	"path"
	"strings"
)

func RepoNameFromURL(rawURL string) (string, error) {
	name := strings.TrimSuffix(path.Base(rawURL), ".git")
	if name == "" || name == "." {
		return "", fmt.Errorf("could not derive repo name from URL: %s", rawURL)
	}
	return name, nil
}
