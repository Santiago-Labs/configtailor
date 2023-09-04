package configtailor

import (
	"path/filepath"
	"strings"
)

func trimPathPrefix(prefix, p string) string {
	// Ensure the prefix has a trailing path separator
	if !strings.HasSuffix(p, string(filepath.Separator)) {
		p += string(filepath.Separator)
	}

	// Use the Rel function to get the relative path
	rel, err := filepath.Rel(prefix, p)
	if err != nil {
		return prefix // Return the original path if there's an error
	}
	return rel
}
