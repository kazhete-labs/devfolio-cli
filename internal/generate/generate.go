package generate

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/kazhetelabs/devfolio-cli/internal/emit"
	"github.com/kazhetelabs/devfolio-cli/internal/github"
	"github.com/kazhetelabs/devfolio-cli/internal/model"
	"github.com/kazhetelabs/devfolio-cli/internal/score"
)

// Options controls a generate run.
type Options struct {
	Login      string
	OutDir     string
	MaxRepos   int
	SkipREADME bool
	Token      string
}

// Result is returned to the CLI for printing.
type Result struct {
	Portfolio model.Portfolio
	OutDir    string
}

// Runner orchestrates fetch → score → emit.
type Runner struct {
	GH *github.Client
}

// NewRunner builds a Runner with a GitHub client.
func NewRunner(token string) *Runner {
	return &Runner{GH: github.New(token)}
}

// Run executes the vertical slice.
func (r *Runner) Run(ctx context.Context, opt Options) (Result, error) {
	if opt.Login == "" {
		return Result{}, fmt.Errorf("github username is required")
	}
	if opt.OutDir == "" {
		opt.OutDir = "devfolio-out"
	}
	if opt.MaxRepos <= 0 {
		opt.MaxRepos = 12
	}

	user, err := r.GH.FetchUser(ctx, opt.Login)
	if err != nil {
		return Result{}, fmt.Errorf("fetch user: %w", err)
	}
	repos, err := r.GH.FetchRepos(ctx, opt.Login)
	if err != nil {
		return Result{}, fmt.Errorf("fetch repos: %w", err)
	}

	// Prefer starred repos first, then name.
	sort.SliceStable(repos, func(i, j int) bool {
		if repos[i].StargazersCount == repos[j].StargazersCount {
			return repos[i].Name < repos[j].Name
		}
		return repos[i].StargazersCount > repos[j].StargazersCount
	})
	if len(repos) > opt.MaxRepos {
		repos = repos[:opt.MaxRepos]
	}

	langs := map[string]int{}
	var sum float64
	scored := 0
	for i := range repos {
		if repos[i].Language != "" {
			langs[repos[i].Language]++
		}
		if !opt.SkipREADME {
			md, err := r.GH.FetchREADME(ctx, opt.Login, repos[i].Name)
			if err != nil {
				return Result{}, fmt.Errorf("readme %s: %w", repos[i].Name, err)
			}
			repos[i].README = md
		}
		repos[i].Score = score.ScoreREADME(repos[i].README)
		if repos[i].Score.Max > 0 {
			sum += float64(repos[i].Score.Total) / float64(repos[i].Score.Max) * 100
			scored++
		}
	}

	avg := 0.0
	if scored > 0 {
		avg = sum / float64(scored)
	}

	p := model.Portfolio{
		User:           user,
		Repos:           repos,
		Languages:      langs,
		AverageScore:   avg,
		GeneratedAtUTC: time.Now().UTC().Format(time.RFC3339),
	}
	if err := emit.WritePortfolio(opt.OutDir, p); err != nil {
		return Result{}, err
	}
	return Result{Portfolio: p, OutDir: opt.OutDir}, nil
}
