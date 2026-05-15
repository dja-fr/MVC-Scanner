package analyzer

import (
	"path/filepath"
	"sort"
	"strings"

	"github.com/pkoukk/tiktoken-go"
)

func getGrade(tokens int) string {
	if tokens < 2000 {
		return "A"
	}
	if tokens < 10000 {
		return "B"
	}
	if tokens < 30000 {
		return "C"
	}
	if tokens < 100000 {
		return "D"
	}
	return "F"
}

// GetGradeIndex returns a numeric value for a grade to allow comparison.
func GetGradeIndex(grade string) int {
	switch grade {
	case "A":
		return 0
	case "B":
		return 1
	case "C":
		return 2
	case "D":
		return 3
	case "F":
		return 4
	}
	return -1
}

// Scan processes all files in the ProjectContext and returns the analysis results.
func Scan(ctx *ProjectContext) ([]FileData, error) {
	tokenizer, err := tiktoken.GetEncoding("cl100k_base")
	if err != nil {
		return nil, err
	}

	var results []FileData

	for _, f := range ctx.Files {
		content := ctx.FileContents[f]

		// Extract dependencies using tree-sitter
		uniqueDeps := ExtractDependencies(f, content, ctx)

		// Tokenization
		selfTokens := len(tokenizer.Encode(string(content), nil, nil))
		totalTokens := selfTokens

		var topDep string
		var topDepTokens int

		for _, d := range uniqueDeps {
			dContent := ctx.FileContents[d]
			dTokens := len(tokenizer.Encode(string(dContent), nil, nil))
			totalTokens += dTokens
			if dTokens > topDepTokens {
				topDepTokens = dTokens
				topDep = d
			}
		}

		relF, _ := filepath.Rel(ctx.AbsRoot, f)

		relTopDep := ""
		if topDep != "" {
			relTopDep, _ = filepath.Rel(ctx.AbsRoot, topDep)
		}

		relDeps := make([]string, len(uniqueDeps))
		for i, d := range uniqueDeps {
			relDeps[i], _ = filepath.Rel(ctx.AbsRoot, d)
		}

		results = append(results, FileData{
			Path:         strings.ReplaceAll(relF, "\\", "/"),
			SelfTokens:   selfTokens,
			TotalTokens:  totalTokens,
			Grade:        getGrade(totalTokens),
			Dependencies: relDeps,
			TopDep:       strings.ReplaceAll(relTopDep, "\\", "/"),
			TopDepTokens: topDepTokens,
		})
	}

	// Aggregate and Sort
	sort.Slice(results, func(i, j int) bool {
		return results[i].TotalTokens > results[j].TotalTokens
	})

	return results, nil
}
