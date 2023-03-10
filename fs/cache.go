package fs

import (
	"os"
	"path/filepath"
)

func CacheDir() (string, error) {
	baseDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return baseDir + filepath.FromSlash("/novasearch"), err
}
