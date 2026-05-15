package analyzer

import (
	"os"
	"path/filepath"
	"strings"
)

// ProjectContext holds the crawled files and their contents.
type ProjectContext struct {
	AbsRoot      string
	Files        []string
	FileContents map[string][]byte
	FileLookup   map[string]string // normalized relative path -> absolute path
}

// Crawl scans the root directory for supported files and loads them into memory.
func Crawl(absRoot string) (*ProjectContext, error) {
	var files []string
	err := filepath.Walk(absRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		ext := filepath.Ext(path)
		if ext == ".java" || ext == ".js" || ext == ".ts" || ext == ".py" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	fileContents := make(map[string][]byte)
	fileLookup := make(map[string]string)

	for _, f := range files {
		content, err := os.ReadFile(f)
		if err == nil {
			fileContents[f] = content
			rel, _ := filepath.Rel(absRoot, f)
			// Normalize separators for lookup
			normalized := strings.ReplaceAll(rel, "\\", "/")
			fileLookup[normalized] = f
		}
	}

	return &ProjectContext{
		AbsRoot:      absRoot,
		Files:        files,
		FileContents: fileContents,
		FileLookup:   fileLookup,
	}, nil
}
