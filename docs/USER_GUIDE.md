# User Guide — devfolio-cli (Go)

Generate a static portfolio site and README quality scorecard from any public GitHub username.

## Requirements

- Go 1.22+ (to build from source), **or** a prebuilt `devfolio` binary
- Network access to `api.github.com`
- Optional: `GITHUB_TOKEN` (classic or fine-grained, public-repo read) for higher rate limits

Without a token, unauthenticated GitHub API limits apply (~60 req/hour). Scoring several repos can hit **403** on README fetches — set a token if that happens.

## Install

```bash
# From module (when published)
go install github.com/kazhetelabs/devfolio-cli/cmd/devfolio@latest

# From this repo
go build -o bin/devfolio ./cmd/devfolio
```

Windows: `go build -o bin/devfolio.exe ./cmd/devfolio`

## Quick start

```bash
# Optional but recommended
set GITHUB_TOKEN=ghp_xxxxxxxx          # PowerShell: $env:GITHUB_TOKEN="..."

devfolio generate YOUR_GITHUB_USER -o ./devfolio-out
```

Then open:

| File | Purpose |
|------|---------|
| `devfolio-out/index.html` | Portfolio page |
| `devfolio-out/scorecard.html` | Browser scorecard |
| `devfolio-out/scorecard.md` | Markdown scorecard (great for PRs/notes) |
| `devfolio-out/styles.css` | Styles for the HTML pages |

## Commands

### `devfolio generate <github-user>`

Fetches the public profile, top non-fork repos (by stars), scores each README, writes the output folder.

| Flag | Default | Description |
|------|---------|-------------|
| `-o`, `--out` | `devfolio-out` | Output directory |
| `--max-repos` | `12` | Max non-fork repos to include |
| `--token` | *(env `GITHUB_TOKEN`)* | GitHub token |
| `--timeout` | `2m` | Overall timeout |

Examples:

```bash
devfolio generate octocat -o ./out --max-repos 5
devfolio generate torvalds --token "$GITHUB_TOKEN"
```

### `devfolio version`

Prints CLI version.

## What the scorecard measures

Each README is scored out of **100**:

| Check | Weight |
|-------|-------:|
| Non-empty | 10 |
| Markdown heading | 10 |
| Length ≥ 400 chars | 10 |
| Install / getting started | 20 |
| Demo / screenshot | 15 |
| Badges | 10 |
| License section | 10 |
| Architecture / how-it-works | 10 |
| Code samples | 5 |

Grades: A (≥90%), B (≥80%), C (≥70%), D (≥55%), F (below).

## Troubleshooting

| Symptom | Fix |
|---------|-----|
| `403 Forbidden` on README | Set `GITHUB_TOKEN`; wait for rate-limit reset |
| `404` user | Check username spelling; profile must be public |
| Empty scorecard rows | Repo has no README — score will be near zero |
| Slow run | Lower `--max-repos`; network dominates runtime |

## Develop / test

```bash
go test ./...
go vet ./...
go build -o bin/devfolio ./cmd/devfolio
```

See also: [README](../README.md), [ADR](adr/0001-architecture.md), [Go vs Rust comparison](../COMPARISON.md).
