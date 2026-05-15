package analyzer

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
)

// ReportOutputs outputs the scan results in the specified format (text, json, csv).
func ReportOutputs(results []FileData, outputFormat string) {
	switch outputFormat {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(results)
	case "csv":
		writer := csv.NewWriter(os.Stdout)
		writer.Write([]string{"Path", "SelfTokens", "TotalTokens", "Grade", "DependenciesCount", "TopContributor", "TopContributorTokens"})
		for _, r := range results {
			writer.Write([]string{
				r.Path,
				fmt.Sprint(r.SelfTokens),
				fmt.Sprint(r.TotalTokens),
				r.Grade,
				fmt.Sprint(len(r.Dependencies)),
				r.TopDep,
				fmt.Sprint(r.TopDepTokens),
			})
		}
		writer.Flush()
	default:
		fmt.Println("========================================================================")
		fmt.Printf("TOKEN DEBT ANALYZER - SCAN RESULTS\n")
		fmt.Printf("Model: cl100k_base | Depth: 1 | Files scanned: %d\n", len(results))
		fmt.Println("========================================================================")
		fmt.Println()

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "GRADE\tFILE\tSELF TOKENS\tTOTAL TOKENS\tDEPS\tTOP CONTRIBUTOR")
		for _, r := range results {
			topDepStr := "-"
			if r.TopDep != "" {
				topDepStr = fmt.Sprintf("%s (%d)", r.TopDep, r.TopDepTokens)
			}
			fmt.Fprintf(w, "%s\t%s\t%d\t%d\t%d\t%s\n", r.Grade, r.Path, r.SelfTokens, r.TotalTokens, len(r.Dependencies), topDepStr)
		}
		w.Flush()
		fmt.Println()
		fmt.Println("========================================================================")
	}
}
