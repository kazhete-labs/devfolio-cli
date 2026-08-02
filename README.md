# devfolio-cli

[![CI](https://github.com/kazhete-labs/devfolio-cli/actions/workflows/ci.yml/badge.svg)](https://github.com/kazhete-labs/devfolio-cli/actions/workflows/ci.yml)
[![Go](https://img.shields.io/badge/Go-1.22%2B-00ADD8?logo=go&logoColor=white)](go.mod)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Rust port](https://img.shields.io/badge/also%20available%20in-Rust-DEA584?logo=rust&logoColor=white)](../devfolio-cli-rust)

Turn any GitHub username into a **static portfolio site** and a **README quality scorecard** — one command, no SaaS, no signup.

```bash
devfolio generate octocat -o ./devfolio-out
```

```text
Generated portfolio for @octocat
  repos scored: 8
  avg README:   42.5 / 100
  output:       ./devfolio-out
  open:         ./devfolio-out/index.html
  scorecard:    ./devfolio-out/scorecard.md
```

## Why devfolio-cli

- 🚀 **Zero config** — point it at a GitHub username, get a finished site
- 📊 **README scorecard** — every repo graded A–F across 9 weighted checks
- 🖥️ **Static output** — plain HTML/CSS, deploy anywhere (GitHub Pages, Netlify, S3…)
- 🔑 **Rate-limit friendly** — optional `GITHUB_TOKEN` for heavier use
- 🐹 **Single binary** — pure Go, no runtime dependencies

## Install

```bash
go install github.com/kazhetelabs/devfolio-cli/cmd/devfolio@latest
```

Or build from this repo:

```bash
go build -o bin/devfolio ./cmd/devfolio
```

## 📖 User Guide

**New here? Start with the [full User Guide](docs/USER_GUIDE.md)** — install options, every flag, the complete scoring rubric, and a troubleshooting table for common errors.

## What it scores

Each README is graded out of 100 — this is the differentiator vs. "another static HTML portfolio generator":

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

## See also

- [User Guide](docs/USER_GUIDE.md) — the full walkthrough
- [Go vs Rust comparison](COMPARISON.md) — benchmarks against the [Rust sibling](../devfolio-cli-rust)

## License

MIT — see [LICENSE](LICENSE).
