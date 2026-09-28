package scanner

import (
	"regexp"
)

type secretPattern struct {
	Name     string
	Regex    *regexp.Regexp
	Severity Severity
}

var secretPatterns = []secretPattern{
	{
		Name: "AWS Access Key",
		Regex: regexp.MustCompile(`\bAKIA[0-9A-Z]{16}\b`),
		Severity: High,
	},
	{
		Name: "GitHub Token",
		Regex: regexp.MustCompile(`\bgh[pousr]_[A-Za-z0-9_]{20,}\b`),
		Severity: High,
	},
	{
		Name: "Generic API Key",
		Regex: regexp.MustCompile(
			`(?i)(api[_-]?key|apikey)\s*[:=]\s*["']?[A-Za-z0-9_\-]{16,}["']?`,
		),
		Severity: High,
	},
	{
		Name: "Generic Password",
		Regex: regexp.MustCompile(
			`(?i)(password|passwd|pwd)\s*[:=]\s*["']?[^"'\s]{8,}["']?`,
		),
		Severity: Medium,
	},
	{
		Name: "Bearer Token",
		Regex: regexp.MustCompile(
			`(?i)bearer\s+[A-Za-z0-9\-._~+/]+=*`,
		),
		Severity: High,
	},
	{
		Name: "Private Key",
		Regex: regexp.MustCompile(
			`-----BEGIN (RSA |EC |OPENSSH )?PRIVATE KEY-----`,
		),
		Severity: High,
	},
}

func DetectSecrets(
	projectName string,
	filePath string,
	content string,
) []Finding {

	var findings []Finding

	for _, pattern := range secretPatterns {
		if pattern.Regex.MatchString(content) {
			findings = append(findings, Finding{
				Project: projectName,
				Issue: pattern.Name +
					" detected in " +
					filePath,
				Severity: pattern.Severity,
			})
		}
	}

	return findings
}
