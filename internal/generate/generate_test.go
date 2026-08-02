package generate

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/kazhetelabs/devfolio-cli/internal/github"
)

func TestRun_EndToEndMockAPI(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/users/octocat", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"login":"octocat","name":"The Octocat","bio":"hi","avatar_url":"https://x/a.png","html_url":"https://github.com/octocat","public_repos":1,"followers":1,"following":0}`))
	})
	mux.HandleFunc("/users/octocat/repos", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"name":"Hello-World","full_name":"octocat/Hello-World","description":"demo","html_url":"https://github.com/octocat/Hello-World","language":"Go","stargazers_count":5,"fork":false}]`))
	})
	mux.HandleFunc("/repos/octocat/Hello-World/readme", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("# Hi\n\n## Install\n\ngo install example.com/x@latest\n"))
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	dir := t.TempDir()
	runner := &Runner{GH: &github.Client{
		HTTP:    srv.Client(),
		BaseURL: srv.URL,
	}}
	res, err := runner.Run(context.Background(), Options{
		Login:    "octocat",
		OutDir:   dir,
		MaxRepos: 5,
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.Portfolio.User.Login != "octocat" {
		t.Fatalf("login=%s", res.Portfolio.User.Login)
	}
	if len(res.Portfolio.Repos) != 1 {
		t.Fatalf("repos=%d", len(res.Portfolio.Repos))
	}
	for _, name := range []string{"index.html", "scorecard.md", "scorecard.html", "styles.css"} {
		if _, err := os.Stat(filepath.Join(dir, name)); err != nil {
			t.Fatalf("missing %s: %v", name, err)
		}
	}
}

func TestRun_RequiresLogin(t *testing.T) {
	runner := NewRunner("")
	_, err := runner.Run(context.Background(), Options{})
	if err == nil {
		t.Fatal("expected error for empty login")
	}
}
