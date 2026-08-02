package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/kazhetelabs/devfolio-cli/internal/model"
)

const defaultBase = "https://api.github.com"

// Client talks to GitHub REST API (public, optional token).
type Client struct {
	HTTP    *http.Client
	BaseURL string
	Token   string
}

// New creates a Client. Token from GITHUB_TOKEN if empty.
func New(token string) *Client {
	if token == "" {
		token = os.Getenv("GITHUB_TOKEN")
	}
	return &Client{
		HTTP:    &http.Client{Timeout: 30 * time.Second},
		BaseURL: defaultBase,
		Token:   token,
	}
}

type apiUser struct {
	Login       string `json:"login"`
	Name        string `json:"name"`
	Bio         string `json:"bio"`
	Company     string `json:"company"`
	Blog        string `json:"blog"`
	Location    string `json:"location"`
	AvatarURL   string `json:"avatar_url"`
	HTMLURL     string `json:"html_url"`
	PublicRepos int    `json:"public_repos"`
	Followers   int    `json:"followers"`
	Following   int    `json:"following"`
}

type apiRepo struct {
	Name            string   `json:"name"`
	FullName        string   `json:"full_name"`
	Description     string   `json:"description"`
	HTMLURL         string   `json:"html_url"`
	Language        string   `json:"language"`
	StargazersCount int      `json:"stargazers_count"`
	ForksCount      int      `json:"forks_count"`
	Fork            bool     `json:"fork"`
	Archived        bool     `json:"archived"`
	Topics          []string `json:"topics"`
	DefaultBranch   string   `json:"default_branch"`
}

// FetchUser loads a public profile.
func (c *Client) FetchUser(ctx context.Context, login string) (model.User, error) {
	var u apiUser
	if err := c.get(ctx, "/users/"+login, &u); err != nil {
		return model.User{}, err
	}
	return model.User{
		Login:       u.Login,
		Name:        u.Name,
		Bio:         u.Bio,
		Company:     u.Company,
		Blog:        u.Blog,
		Location:    u.Location,
		AvatarURL:   u.AvatarURL,
		HTMLURL:     u.HTMLURL,
		PublicRepos: u.PublicRepos,
		Followers:   u.Followers,
		Following:   u.Following,
	}, nil
}

// FetchRepos lists public non-fork repos (up to 100), sorted by updated.
func (c *Client) FetchRepos(ctx context.Context, login string) ([]model.Repo, error) {
	var raw []apiRepo
	path := fmt.Sprintf("/users/%s/repos?per_page=100&sort=updated&type=owner", login)
	if err := c.get(ctx, path, &raw); err != nil {
		return nil, err
	}
	out := make([]model.Repo, 0, len(raw))
	for _, r := range raw {
		if r.Fork {
			continue
		}
		out = append(out, model.Repo{
			Name:            r.Name,
			FullName:        r.FullName,
			Description:     r.Description,
			HTMLURL:         r.HTMLURL,
			Language:        r.Language,
			StargazersCount: r.StargazersCount,
			ForksCount:      r.ForksCount,
			Fork:            r.Fork,
			Archived:        r.Archived,
			Topics:          r.Topics,
			DefaultBranch:   r.DefaultBranch,
		})
	}
	return out, nil
}

// FetchREADME downloads README markdown for a repo (empty string if none).
func (c *Client) FetchREADME(ctx context.Context, owner, repo string) (string, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/readme", c.BaseURL, owner, repo)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	c.applyHeaders(req)
	req.Header.Set("Accept", "application/vnd.github.raw")

	res, err := c.HTTP.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusNotFound {
		return "", nil
	}
	if res.StatusCode >= 400 {
		body, _ := io.ReadAll(io.LimitReader(res.Body, 2048))
		return "", fmt.Errorf("github README %s: %s", res.Status, strings.TrimSpace(string(body)))
	}
	b, err := io.ReadAll(io.LimitReader(res.Body, 1<<20))
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (c *Client) get(ctx context.Context, path string, dest any) error {
	url := c.BaseURL + path
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	c.applyHeaders(req)
	req.Header.Set("Accept", "application/vnd.github+json")

	res, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode >= 400 {
		body, _ := io.ReadAll(io.LimitReader(res.Body, 2048))
		return fmt.Errorf("github %s: %s", res.Status, strings.TrimSpace(string(body)))
	}
	return json.NewDecoder(res.Body).Decode(dest)
}

func (c *Client) applyHeaders(req *http.Request) {
	req.Header.Set("User-Agent", "devfolio-cli")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}
}
