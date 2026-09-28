package scanner

import (
	"fmt"
	"strings"
	"sync"

	"gitlab-risk-scanner/internal/gitlab"
)

type Severity string

const (
	High   Severity = "HIGH"
	Medium Severity = "MEDIUM"
	Low    Severity = "LOW"
)

type Finding struct {
	Project  string
	Issue    string
	Severity Severity
}

type Scanner struct {
	Client      *gitlab.Client
	Concurrency int
}

func New(client *gitlab.Client, concurrency int) *Scanner {
	if concurrency <= 0 {
		concurrency = 5
	}

	return &Scanner{
		Client:      client,
		Concurrency: concurrency,
	}
}

func (s *Scanner) ScanProjects(projects []gitlab.Project) []Finding {
	var (
		wg       sync.WaitGroup
		mu       sync.Mutex
		findings []Finding
	)

	jobs := make(chan gitlab.Project)

	for i := 0; i < s.Concurrency; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for project := range jobs {
				projectFindings := s.scanProject(project)

				mu.Lock()
				findings = append(findings, projectFindings...)
				mu.Unlock()
			}
		}()
	}

	for _, project := range projects {
		jobs <- project
	}

	close(jobs)

	wg.Wait()

	return findings
}

func (s *Scanner) scanProject(project gitlab.Project) []Finding {
	var findings []Finding

	tree, err := s.Client.GetRepositoryTree(project.ID)
	if err != nil {
		return []Finding{
			{
				Project:  project.Name,
				Issue:    fmt.Sprintf("Unable to scan repository: %v", err),
				Severity: Medium,
			},
		}
	}

	findings = append(
		findings,
		CheckSensitiveFiles(project, tree)...,
	)

	findings = append(
		findings,
		CheckMetadata(project, tree)...,
	)

	for _, item := range tree {
		if item.Type != "blob" {
			continue
		}

		// Avoid scanning obviously binary/private-key files as text
		// unless the secret detector specifically handles them.
		if !shouldScanContent(item.Path) {
			continue
		}

		content, err := s.Client.GetFile(
			project.ID,
			item.Path,
			project.DefaultBranch,
		)

		if err != nil {
			continue
		}

		secretFindings := DetectSecrets(
			project.Name,
			item.Path,
			content,
		)

		findings = append(findings, secretFindings...)
	}

	return findings
}

func shouldScanContent(path string) bool {
	lower := strings.ToLower(path)

	binaryExtensions := []string{
		".png",
		".jpg",
		".jpeg",
		".gif",
		".pdf",
		".zip",
		".tar",
		".gz",
		".exe",
		".dll",
		".so",
		".woff",
		".woff2",
	}

	for _, ext := range binaryExtensions {
		if strings.HasSuffix(lower, ext) {
			return false
		}
	}

	return true
}
