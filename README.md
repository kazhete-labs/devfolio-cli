# devfolio-cli

Turn a GitHub username into a **static portfolio site** plus a **README quality scorecard**.

Local CLI only — no SaaS. Works with the public GitHub API (optional `GITHUB_TOKEN` for higher rate limits).

## User guide

Full install, flags, scoring rules, and troubleshooting: **[docs/USER_GUIDE.md](docs/USER_GUIDE.md)**.

## Install

```bash
go install github.com/kazhetelabs/devfolio-cli/cmd/devfolio@latest
```

Or from this repo:

```bash
go build -o bin/devfolio ./cmd/devfolio
```

## 30-second demo

```bash
devfolio generate octocat -o ./devfolio-out
# open ./devfolio-out/index.html
# read ./devfolio-out/scorecard.md
```

Example terminal output:

```text
Generated portfolio for @octocat
  repos scored: 8
  avg README:   42.5 / 100
  output:       ./devfolio-out
  open:         ./devfolio-out/index.html
  scorecard:    ./devfolio-out/scorecard.md
```

> Tip: add a short screen recording as `docs/demo.gif` before publishing the repo publicly.

## What it scores

Each README is graded out of 100 on:

| Check | Weight |
|-------|-------:|
| Non-empty | 10 |
| Heading | 10 |
| Length ≥ 400 chars | 10 |
| Install / getting started | 20 |
| Demo / screenshot | 15 |
| Badges | 10 |
| License section | 10 |
| Architecture / how-it-works | 10 |
| Code samples | 5 |

This is the differentiator vs “another static HTML portfolio generator.”

## Architecture

```text
┌─────────────┐     ┌──────────────────┐     ┌─────────────────┐
│ cobra CLI   │───▶│ generate runner  │───▶│ emit (HTML/MD)  │
│ (delivery)  │     │ (application)    │     │ (adapter)       │
└─────────────┘     └────────┬─────────┘     └─────────────────┘
                             │
              ┌──────────────┼──────────────┐
              ▼                             ▼
     ┌─────────────────┐          ┌─────────────────┐
     │ github REST     │          │ score.README    │
     │ (adapter)       │          │ (domain)        │
     └─────────────────┘          └─────────────────┘
```

ADR: [`docs/adr/0001-architecture.md`](docs/adr/0001-architecture.md)

## Flags

```text
devfolio generate <user> [flags]
  -o, --out string         output directory (default "devfolio-out")
      --max-repos int      max non-fork repos (default 12)
      --token string       GitHub token (or GITHUB_TOKEN)
      --timeout duration   overall timeout (default 2m)
```

## Development

```bash
go test ./...
go vet ./...
go build -o bin/devfolio ./cmd/devfolio
```

## License

MIT — see [LICENSE](LICENSE).
