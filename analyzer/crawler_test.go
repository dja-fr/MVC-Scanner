package analyzer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCrawl(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "crawler_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir) // Cleanup after test

	// Create some dummy files
	filesToCreate := []string{
		"test.java",
		"index.js",
		"main.ts",
		"script.py",
		"ignore.txt", // Should be ignored
	}

	for _, f := range filesToCreate {
		path := filepath.Join(tempDir, f)
		err := os.WriteFile(path, []byte("dummy content"), 0644)
		if err != nil {
			t.Fatalf("Failed to write dummy file: %v", err)
		}
	}

	// Also create a sub-directory and a file inside
	subDir := filepath.Join(tempDir, "src")
	err = os.Mkdir(subDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create sub dir: %v", err)
	}
	err = os.WriteFile(filepath.Join(subDir, "service.ts"), []byte("dummy content"), 0644)
	if err != nil {
		t.Fatalf("Failed to write file in sub dir: %v", err)
	}

	// Run crawler
	ctx, err := Crawl(tempDir)
	if err != nil {
		t.Fatalf("Crawl failed: %v", err)
	}

	// Assertions
	expectedFilesCount := 5 // java, js, ts, py, and src/service.ts
	if len(ctx.Files) != expectedFilesCount {
		t.Errorf("Expected %d files, got %d", expectedFilesCount, len(ctx.Files))
	}

	if len(ctx.FileContents) != expectedFilesCount {
		t.Errorf("Expected %d file contents, got %d", expectedFilesCount, len(ctx.FileContents))
	}

	if len(ctx.FileLookup) != expectedFilesCount {
		t.Errorf("Expected %d file lookups, got %d", expectedFilesCount, len(ctx.FileLookup))
	}

	// Ensure ignore.txt is not in there
	for _, f := range ctx.Files {
		if filepath.Base(f) == "ignore.txt" {
			t.Errorf("Crawl should have ignored 'ignore.txt' but it was included")
		}
	}
}
