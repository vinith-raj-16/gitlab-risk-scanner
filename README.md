# GitLab Risk Scanner

A command-line tool that scans GitLab repositories for common security risks and basic project hygiene issues. It uses the GitLab REST API and reports findings in a terminal table.

## Requirements

- Go 1.23.4 or newer
- A GitLab personal access token when scanning projects that require authentication

## Build

```powershell
go build -o gitlab-risk-scanner.exe ./cmd/scanner
```

Alternatively, run directly with `go run ./cmd/scanner`.

## Usage

Choose exactly one target: a GitLab user, a group, or a project path.

```powershell
# Scan public repositories owned by a user
go run ./cmd/scanner --user example-user

# Scan public repositories in a group, including subgroups
go run ./cmd/scanner --group example-group

# Scan one project (use the full namespace/project path)
go run ./cmd/scanner --project example-group/example-project

# Supply a token and scan up to 8 repositories concurrently
go run ./cmd/scanner --group example-group --token $env:GITLAB_TOKEN --concurrency 8
```

`--token` is optional for public projects. Set `GITLAB_TOKEN` in your PowerShell session before passing it on the command line:

```powershell
$env:GITLAB_TOKEN = "your-personal-access-token"
```

Use a token with access to the project when scanning private repositories. Avoid committing tokens or placing them in shared scripts.

## Checks

- **Secrets:** common AWS access keys, GitHub tokens, API keys, passwords, bearer tokens, and private-key markers in repository file contents.
- **Sensitive files:** committed `.env`, SSH private-key files, `secrets.yml`/`secrets.yaml`, PEM files, and `config.json`.
- **Project metadata:** missing README and LICENSE files.

Secret checks use pattern matching and can produce false positives or miss unsupported formats. This tool is a lightweight review aid, not a replacement for a dedicated secret scanner or security audit. Repository content is checked through the GitLab API; the scanner does not clone projects.
