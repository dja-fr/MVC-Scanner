package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"token-debt-analyzer/analyzer"
)

var (
	outputFormat = flag.String("output", "text", "Output format: text, json, csv")
	failOn       = flag.String("fail-on", "", "Fail CI if at least one file reaches this grade or worse (e.g., D)")
)

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) < 2 || args[0] != "scan" {
		fmt.Println("Usage: token-debt-analyzer scan <path>")
		os.Exit(1)
	}
	rootDir := args[1]
	absRoot, err := filepath.Abs(rootDir)
	if err != nil {
		fmt.Printf("Error resolving path: %v\n", err)
		os.Exit(1)
	}

	// 1. Crawl
	ctx, err := analyzer.Crawl(absRoot)
	if err != nil {
		fmt.Printf("Error crawling directory: %v\n", err)
		os.Exit(1)
	}

	// 2 & 3 & 4. Parse, Resolve, Tokenize and Aggregate
	results, err := analyzer.Scan(ctx)
	if err != nil {
		fmt.Printf("Error during scanning: %v\n", err)
		os.Exit(1)
	}

	// 5. Output
	analyzer.ReportOutputs(results, *outputFormat)

	if *failOn != "" {
		failIndex := analyzer.GetGradeIndex(*failOn)
		if failIndex != -1 {
			failed := false
			for _, r := range results {
				if analyzer.GetGradeIndex(r.Grade) >= failIndex {
					failed = true
					break
				}
			}
			if failed {
				os.Exit(1)
			}
		}
	}
}
