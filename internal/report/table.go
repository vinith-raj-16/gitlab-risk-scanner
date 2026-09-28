package report

import (
	"fmt"
	"gitlab-risk-scanner/internal/scanner"
	"sort"
)

func PrintTable(findings []scanner.Finding) {
	if len(findings) == 0 {
		fmt.Println("\nNo risks detected.")
		return
	}

	sort.Slice(findings, func(i, j int) bool {
		if findings[i].Project == findings[j].Project {
			return findings[i].Severity < findings[j].Severity
		}

		return findings[i].Project < findings[j].Project
	})

	fmt.Println()

	fmt.Printf(
		"%-30s %-55s %-10s\n",
		"PROJECT",
		"ISSUE",
		"SEVERITY",
	)

	fmt.Println("--------------------------------------------------------------------------------------------------------")

	for _, finding := range findings {
		fmt.Printf(
			"%-30s %-55s %-10s\n",
			truncate(finding.Project, 30),
			truncate(finding.Issue, 55),
			finding.Severity,
		)
	}

	fmt.Println()
	fmt.Printf("Total findings: %d\n", len(findings))
}

func truncate(value string, max int) string {
	if len(value) <= max {
		return value
	}

	return value[:max-3] + "..."
}
