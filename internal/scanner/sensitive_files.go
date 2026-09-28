package scanner

import (
	"path/filepath"
	"strings"

	"gitlab-risk-scanner/internal/gitlab"
)

var sensitiveFileNames = map[string]Severity{
	".env":         High,
	"id_rsa":       High,
	"id_dsa":       High,
	"id_ecdsa":     High,
	"id_ed25519":   High,
	"secrets.yml":  High,
	"secrets.yaml": High,
	".pem":         High,
	"config.json":  Medium,
}

func CheckSensitiveFiles(
	project gitlab.Project,
	tree []gitlab.TreeItem,
) []Finding {

	var findings []Finding

	for _, item := range tree {
		if item.Type != "blob" {
			continue
		}

		base := strings.ToLower(filepath.Base(item.Path))

		if severity, exists := sensitiveFileNames[base]; exists {
			findings = append(findings, Finding{
				Project: project.Name,
				Issue:    "Sensitive file committed: " + item.Path,
				Severity: severity,
			})

			continue
		}

		// Check extensions separately.
		if strings.HasSuffix(base, ".pem") {
			findings = append(findings, Finding{
				Project: project.Name,
				Issue:    "Private/certificate file committed: " + item.Path,
				Severity: High,
			})
		}
	}

	return findings
}
