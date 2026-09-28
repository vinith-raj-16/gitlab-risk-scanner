package scanner

import (
	"strings"

//	"gitlab-risk-scanner/internal/gitlab"
"gitlab-risk-scanner/internal/gitlab"
)

func CheckMetadata(
	project gitlab.Project,
	tree []gitlab.TreeItem,
) []Finding {

	var (
		findings    []Finding
		hasReadme   bool
		hasLicense  bool
	)

	for _, item := range tree {
		if item.Type != "blob" {
			continue
		}

		name := strings.ToLower(item.Name)

		switch name {
		case "readme.md", "readme":
			hasReadme = true

		case "license", "license.md", "license.txt":
			hasLicense = true
		}
	}

	if !hasReadme {
		findings = append(findings, Finding{
			Project:  project.Name,
			Issue:    "README.md is missing",
			Severity: Low,
		})
	}

	if !hasLicense {
		findings = append(findings, Finding{
			Project:  project.Name,
			Issue:    "LICENSE is missing",
			Severity: Low,
		})
	}

	return findings
}
