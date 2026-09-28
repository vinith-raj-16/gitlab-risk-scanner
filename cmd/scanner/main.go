package main

import (
	"flag"
	"fmt"
	"os"

	"gitlab-risk-scanner/internal/gitlab"
	"gitlab-risk-scanner/internal/report"
	"gitlab-risk-scanner/internal/scanner"
)

func main() {
	username := flag.String(
		"user",
		"",
		"GitLab username",
	)

	group := flag.String(
		"group",
		"",
		"GitLab group",
	)

	project := flag.String(
		"project",
		"",
		"GitLab project path, e.g. username/repository",
	)

	token := flag.String(
		"token",
		"",
		"GitLab personal access token",
	)

	concurrency := flag.Int(
		"concurrency",
		5,
		"Number of repositories to scan concurrently",
	)

	flag.Parse()

	if *username == "" && *group == "" && *project == "" {
		fmt.Println("Error: provide --user, --group, or --project")
		flag.Usage()
		os.Exit(1)
	}

	targetCount := 0

	if *username != "" {
		targetCount++
	}

	if *group != "" {
		targetCount++
	}

	if *project != "" {
		targetCount++
	}

	if targetCount > 1 {
		fmt.Println(
			"Error: use only one of --user, --group, or --project",
		)
		os.Exit(1)
	}

	client := gitlab.NewClient(*token)

	var (
		projects []gitlab.Project
		err      error
	)

	if *project != "" {

		fmt.Printf(
			"Getting project: %s\n",
			*project,
		)

		var p gitlab.Project

		p, err = client.GetProject(*project)

		if err == nil {
			projects = []gitlab.Project{p}
		}

	} else if *username != "" {

		fmt.Printf(
			"Finding public repositories for user: %s\n",
			*username,
		)

		projects, err = client.GetUserProjects(*username)

	} else if *group != "" {

		fmt.Printf(
			"Finding public repositories for group: %s\n",
			*group,
		)

		projects, err = client.GetGroupProjects(*group)
	}

	if err != nil {
		fmt.Printf(
			"Failed to get projects: %v\n",
			err,
		)
		os.Exit(1)
	}

	if len(projects) == 0 {
		fmt.Println("No public repositories found.")
		return
	}

	fmt.Printf(
		"Found %d repositories.\n",
		len(projects),
	)

	s := scanner.New(
		client,
		*concurrency,
	)

	findings := s.ScanProjects(projects)

	report.PrintTable(findings)
}
