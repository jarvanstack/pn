package utils

import (
	"os"
	"path"
	"path/filepath"
)

func MatchFiles(root string, pattern string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match(pattern, info.Name()); err != nil {
			return err
		} else if matched {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

func RandomFileName(fileName string) string {
	suffix := path.Ext(fileName)
	return UUID() + suffix
}
