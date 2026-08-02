package github

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchUserAndRepos(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/users/octocat", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"login":"octocat","name":"The Octocat","bio":"hi","avatar_url":"https://x/a.png","html_url":"https://github.com/octocat","public_repos":2,"followers":1,"following":1}`))
	})
	mux.HandleFunc("/users/octocat/repos", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"name":"Hello-World","full_name":"octocat/Hello-World","description":"demo","html_url":"https://github.com/octocat/Hello-World","language":"Go","stargazers_count":10,"forks_count":1,"fork":false,"archived":false,"topics":["demo"],"default_branch":"main"},{"name":"forked","full_name":"octocat/forked","fork":true,"html_url":"https://github.com/octocat/forked"}]`))
	})
	mux.HandleFunc("/repos/octocat/Hello-World/readme", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte("# Hello\n\n## Install\n\ngo install example.com/x@latest\n"))
	})

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	c := New("")
	c.BaseURL = srv.URL
	c.HTTP = srv.Client()

	u, err := c.FetchUser(context.Background(), "octocat")
	if err != nil {
		t.Fatal(err)
	}
	if u.Login != "octocat" {
		t.Fatalf("login=%s", u.Login)
	}
	repos, err := c.FetchRepos(context.Background(), "octocat")
	if err != nil {
		t.Fatal(err)
	}
	if len(repos) != 1 {
		t.Fatalf("expected 1 non-fork repo, got %d", len(repos))
	}
	md, err := c.FetchREADME(context.Background(), "octocat", "Hello-World")
	if err != nil {
		t.Fatal(err)
	}
	if md == "" {
		t.Fatal("expected readme body")
	}
}
