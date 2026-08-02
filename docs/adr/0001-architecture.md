# ADR 0001: Modular CLI with ports & adapters

## Status

Accepted — 2026-08-02

## Context

`devfolio-cli` is a solo-dev, local-first tool: fetch a public GitHub profile, score READMEs, emit a static portfolio + scorecard. No SaaS, no database, free API only.

## Decision

Ship a **single Go binary** structured as a small modular monolith with hexagonal boundaries:

| Layer | Package | Responsibility |
|-------|---------|----------------|
| Delivery | `cmd/devfolio`, `internal/cli` | Cobra commands, flags, UX |
| Domain | `internal/score`, `internal/model` | README scoring rules, portfolio model |
| Application | `internal/generate` | Orchestrate fetch → score → emit |
| Adapters | `internal/github`, `internal/emit` | GitHub REST client, static HTML/MD writers |

No microservices, no shared DB. Optional `GITHUB_TOKEN` for higher rate limits.

## Alternatives considered

1. **Rust + clap** — strong binary story; slower solo iteration here (Go already chosen in Odoo task).
2. **Node CLI** — easier HTML templating; weaker “install one binary” signal for this portfolio goal.
3. **Hosted SaaS** — out of scope; rejected by product constraints.

## Consequences

- Easy to test domain scoring without network (fixture READMEs).
- GitHub client is swappable (REST now; GraphQL pinned-repos later).
- Static emit keeps demos offline-friendly after one fetch.
