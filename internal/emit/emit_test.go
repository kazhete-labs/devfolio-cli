package emit

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kazhetelabs/devfolio-cli/internal/model"
)

func TestWritePortfolio(t *testing.T) {
	dir := t.TempDir()
	p := model.Portfolio{
		User: model.User{
			Login:     "octocat",
			Name:      "The Octocat",
			Bio:       "GitHub mascot",
			AvatarURL: "https://example.com/a.png",
			HTMLURL:   "https://github.com/octocat",
		},
		Repos: []model.Repo{{
			Name:            "Hello-World",
			Description:     "demo",
			HTMLURL:         "https://github.com/octocat/Hello-World",
			Language:        "Go",
			StargazersCount: 3,
			Score: model.READMEScore{
				Total: 50, Max: 100, Grade: "D", Summary: "Improve: install",
				Checks: []model.CheckResult{{ID: "install", Label: "Install", Passed: false, Weight: 20, Detail: "missing"}},
			},
		}},
		Languages:      map[string]int{"Go": 1},
		AverageScore:   50,
		GeneratedAtUTC: "2026-08-02T00:00:00Z",
	}
	if err := WritePortfolio(dir, p); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"index.html", "scorecard.html", "scorecard.md", "styles.css"} {
		if _, err := os.Stat(filepath.Join(dir, name)); err != nil {
			t.Fatalf("missing %s: %v", name, err)
		}
	}
	md, err := os.ReadFile(filepath.Join(dir, "scorecard.md"))
	if err != nil {
		t.Fatal(err)
	}
	if len(md) < 20 {
		t.Fatal("scorecard.md too short")
	}
}
