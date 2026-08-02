package emit

import (
	"fmt"
	"html"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/kazhetelabs/devfolio-cli/internal/model"
)

// WritePortfolio writes index.html, scorecard.md, scorecard.html, and styles.css.
func WritePortfolio(outDir string, p model.Portfolio) error {
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}
	files := map[string]string{
		"styles.css":     stylesCSS,
		"index.html":     renderIndex(p),
		"scorecard.html": renderScorecardHTML(p),
		"scorecard.md":   renderScorecardMD(p),
	}
	for name, body := range files {
		path := filepath.Join(outDir, name)
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			return fmt.Errorf("write %s: %w", name, err)
		}
	}
	return nil
}

func renderIndex(p model.Portfolio) string {
	name := p.User.Name
	if name == "" {
		name = p.User.Login
	}
	var repos strings.Builder
	for _, r := range p.Repos {
		lang := r.Language
		if lang == "" {
			lang = "n/a"
		}
		repos.WriteString(fmt.Sprintf(`
      <article class="card">
        <h3><a href="%s">%s</a></h3>
        <p class="muted">%s</p>
        <p class="meta">★ %d · %s · README %s (%d/%d)</p>
      </article>`,
			html.EscapeString(r.HTMLURL),
			html.EscapeString(r.Name),
			html.EscapeString(r.Description),
			r.StargazersCount,
			html.EscapeString(lang),
			html.EscapeString(r.Score.Grade),
			r.Score.Total,
			r.Score.Max,
		))
	}

	langs := sortedLangLines(p.Languages)
	return fmt.Sprintf(`<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8"/>
  <meta name="viewport" content="width=device-width, initial-scale=1"/>
  <title>%s · Devfolio</title>
  <link rel="stylesheet" href="styles.css"/>
</head>
<body>
  <header class="hero">
    <img class="avatar" src="%s" alt="avatar"/>
    <div>
      <p class="eyebrow">devfolio-cli</p>
      <h1>%s</h1>
      <p class="bio">%s</p>
      <p class="meta"><a href="%s">@%s</a> · %d public repos · %d followers · avg README %.0f/100</p>
    </div>
  </header>
  <main>
    <section>
      <h2>Languages</h2>
      <ul class="langs">%s</ul>
    </section>
    <section>
      <h2>Repositories</h2>
      <div class="grid">%s</div>
    </section>
    <p class="footer"><a href="scorecard.html">Open README scorecard</a> · generated %s UTC</p>
  </main>
</body>
</html>
`,
		html.EscapeString(name),
		html.EscapeString(p.User.AvatarURL),
		html.EscapeString(name),
		html.EscapeString(p.User.Bio),
		html.EscapeString(p.User.HTMLURL),
		html.EscapeString(p.User.Login),
		p.User.PublicRepos,
		p.User.Followers,
		p.AverageScore,
		langs,
		repos.String(),
		html.EscapeString(p.GeneratedAtUTC),
	)
}

func renderScorecardHTML(p model.Portfolio) string {
	var body strings.Builder
	for _, r := range p.Repos {
		body.WriteString(fmt.Sprintf("<h3>%s — %s (%d/%d)</h3><ul>", html.EscapeString(r.Name), html.EscapeString(r.Score.Grade), r.Score.Total, r.Score.Max))
		for _, c := range r.Score.Checks {
			mark := "FAIL"
			if c.Passed {
				mark = "PASS"
			}
			body.WriteString(fmt.Sprintf("<li><strong>%s</strong> [%s] %s — %s</li>", html.EscapeString(mark), html.EscapeString(c.ID), html.EscapeString(c.Label), html.EscapeString(c.Detail)))
		}
		body.WriteString("</ul>")
	}
	return fmt.Sprintf(`<!doctype html>
<html lang="en"><head><meta charset="utf-8"/><title>Scorecard · %s</title><link rel="stylesheet" href="styles.css"/></head>
<body><main><p><a href="index.html">← Portfolio</a></p><h1>README scorecard</h1><p class="muted">@%s · avg %.0f/100</p>%s</main></body></html>`,
		html.EscapeString(p.User.Login), html.EscapeString(p.User.Login), p.AverageScore, body.String())
}

func renderScorecardMD(p model.Portfolio) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("# README scorecard — @%s\n\n", p.User.Login))
	b.WriteString(fmt.Sprintf("Average score: **%.1f / 100** · generated %s UTC\n\n", p.AverageScore, p.GeneratedAtUTC))
	for _, r := range p.Repos {
		b.WriteString(fmt.Sprintf("## %s — grade %s (%d/%d)\n\n", r.Name, r.Score.Grade, r.Score.Total, r.Score.Max))
		b.WriteString(fmt.Sprintf("%s\n\n", r.Score.Summary))
		b.WriteString("| Check | Result | Weight | Detail |\n|---|---|---:|---|\n")
		for _, c := range r.Score.Checks {
			res := "FAIL"
			if c.Passed {
				res = "PASS"
			}
			b.WriteString(fmt.Sprintf("| %s | %s | %d | %s |\n", c.Label, res, c.Weight, c.Detail))
		}
		b.WriteString("\n")
	}
	return b.String()
}

func sortedLangLines(langs map[string]int) string {
	type kv struct {
		k string
		v int
	}
	items := make([]kv, 0, len(langs))
	for k, v := range langs {
		items = append(items, kv{k, v})
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].v == items[j].v {
			return items[i].k < items[j].k
		}
		return items[i].v > items[j].v
	})
	var b strings.Builder
	for _, it := range items {
		b.WriteString(fmt.Sprintf("<li><span>%s</span><span>%d</span></li>", html.EscapeString(it.k), it.v))
	}
	if b.Len() == 0 {
		return "<li><span>n/a</span><span>0</span></li>"
	}
	return b.String()
}

const stylesCSS = `:root {
  --bg: #f4f1ea;
  --ink: #1c1917;
  --muted: #57534e;
  --card: #fffdf8;
  --line: #d6d3d1;
  --accent: #0f766e;
  --font: "Iowan Old Style", "Palatino Linotype", Palatino, Georgia, serif;
  --sans: "Segoe UI", system-ui, sans-serif;
}
* { box-sizing: border-box; }
body {
  margin: 0;
  color: var(--ink);
  background:
    radial-gradient(circle at 10% 10%, #fde68a55, transparent 40%),
    radial-gradient(circle at 90% 0%, #99f6e455, transparent 35%),
    var(--bg);
  font-family: var(--sans);
  line-height: 1.5;
}
.hero, main { max-width: 960px; margin: 0 auto; padding: 2rem 1.25rem; }
.hero { display: flex; gap: 1.25rem; align-items: center; }
.avatar { width: 96px; height: 96px; border-radius: 20px; object-fit: cover; border: 1px solid var(--line); }
h1, h2, h3 { font-family: var(--font); letter-spacing: -0.02em; }
.eyebrow { text-transform: uppercase; letter-spacing: 0.12em; font-size: 0.75rem; color: var(--accent); margin: 0 0 0.25rem; }
.bio { color: var(--muted); max-width: 42rem; }
.meta, .muted, .footer { color: var(--muted); font-size: 0.95rem; }
.grid { display: grid; gap: 1rem; grid-template-columns: repeat(auto-fit, minmax(240px, 1fr)); }
.card { background: var(--card); border: 1px solid var(--line); border-radius: 14px; padding: 1rem 1.1rem; }
.card a { color: var(--ink); text-decoration: none; }
.card a:hover { color: var(--accent); }
.langs { list-style: none; padding: 0; display: grid; gap: 0.35rem; max-width: 320px; }
.langs li { display: flex; justify-content: space-between; border-bottom: 1px dashed var(--line); padding: 0.25rem 0; }
a { color: var(--accent); }
`
