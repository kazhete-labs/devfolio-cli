# Go vs Rust — QA results

Measured on Windows (local), 2026-08-02. Same feature set: `generate <user>` → portfolio HTML + README scorecard.

## Test pyramid (both)

| Layer | Coverage |
|-------|----------|
| Unit | README scoring (table + edge cases) |
| Integration | GitHub client against mock HTTP |
| Vertical slice | generate → emit artifacts (mock API) |

**Result:** Go `go test ./...` **PASS**. Rust `cargo test` **PASS** (7 unit/integration tests).

## Metrics

| Metric | Go | Rust | Winner |
|--------|---:|-----:|--------|
| Release binary size | **6.77 MB** (`-ldflags=-s -w`) | **5.64 MB** (`cargo build --release`) | Rust (~17% smaller) |
| README score ×20k | **0.37 s** | **0.043 s** | Rust (~8.5×) |
| `generate octocat --max-repos 5` (avg of 3, network) | **~458 ms** | **~476 ms** | Tie (API-bound) |
| Cold `--version` / `version` | ~33 ms | ~17 ms | Rust |
| Incremental test run | ~2.0 s | ~0.3 s (cached) | Rust* |
| Clean compile feel | Fast | Slower first build | Go |

\*Rust test time after deps cached; first `cargo test` downloads/compiles deps (~20–40s).

## Qualitative

| Dimension | Go | Rust |
|-----------|----|------|
| DX / iterate | Faster edit-run loop | Stronger type guarantees, longer compile |
| Deps | Small (cobra) | Larger tree (reqwest/rustls) but static TLS |
| Distribution | `go install` simple | `cargo install` / release binary |
| Portfolio signal | Very clear for backend roles | Strong systems/CLI signal |

## Verdict

- **Ship Rust** if you optimize for **binary size + CPU-bound scoring** and want a systems-language pin.
- **Ship Go** if you optimize for **fast iteration** and a slightly simpler contribution story.
- For this CLI’s real workload (`generate`), **network dominates** — runtime difference is negligible end-to-end.

Recommendation for GitHub boost: **pin Rust as primary** (smaller binary, faster scoring) and keep Go as a sibling “same API, two languages” flex — or vice versa if your target employers are Go-heavy.
