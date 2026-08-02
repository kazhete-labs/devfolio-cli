package score

import (
	"regexp"
	"strings"

	"github.com/kazhetelabs/devfolio-cli/internal/model"
)

var (
	reHeading     = regexp.MustCompile(`(?m)^#{1,3}\s+\S+`)
	reBadge       = regexp.MustCompile(`(?i)!\[.*?\]\(https?://.*(shields\.io|badge|img\.shields\.io).*\)|\[!\[.*?\]\(https?://.*\)\]\(https?://.*\)`)
	reInstall     = regexp.MustCompile(`(?im)^#{1,3}\s+.*(install|getting started|quick ?start|usage).*$|go install |npm (i|install) |pip install |cargo install |docker (compose )?up`)
	reDemo        = regexp.MustCompile(`(?im)^#{1,3}\s+.*(demo|screenshot|preview).*$|\.(gif|webm|mp4)\)|ASCIinema|loom\.com`)
	reLicense     = regexp.MustCompile(`(?im)^#{1,3}\s+.*license.*$|mit license|apache license|spdx-license`)
	reArch        = regexp.MustCompile(`(?im)^#{1,3}\s+.*(architecture|design|how it works|overview).*$`)
	reCodeFence   = regexp.MustCompile("(?m)^```")
	reContributing = regexp.MustCompile(`(?im)^#{1,3}\s+.*contribut`)
)

// ScoreREADME evaluates README markdown quality (max 100).
func ScoreREADME(md string) model.READMEScore {
	checks := []model.CheckResult{
		check("non_empty", "README is non-empty", 10, strings.TrimSpace(md) != "", "file missing or empty"),
		check("title", "Has a markdown heading", 10, reHeading.MatchString(md), "add an H1/H2 title"),
		check("length", "Adequate length (≥400 chars)", 10, len(strings.TrimSpace(md)) >= 400, "expand the README"),
		check("install", "Install / getting started section", 20, reInstall.MatchString(md), "document how to install/run"),
		check("demo", "Demo / screenshot evidence", 15, reDemo.MatchString(md), "add a GIF, screenshot, or demo link"),
		check("badges", "Status badges", 10, reBadge.MatchString(md), "add shields.io (or similar) badges"),
		check("license", "License mentioned", 10, reLicense.MatchString(md), "add a License section"),
		check("architecture", "Architecture / how-it-works", 10, reArch.MatchString(md), "sketch architecture briefly"),
		check("code_samples", "Code fences / examples", 5, reCodeFence.MatchString(md), "include copy-pasteable examples"),
	}

	total, max := 0, 0
	for _, c := range checks {
		max += c.Weight
		if c.Passed {
			total += c.Weight
		}
	}

	return model.READMEScore{
		Total:   total,
		Max:     max,
		Checks:  checks,
		Grade:   grade(total, max),
		Summary: summarize(total, max, checks),
	}
}

func check(id, label string, weight int, passed bool, detail string) model.CheckResult {
	if passed {
		detail = "ok"
	}
	return model.CheckResult{ID: id, Label: label, Passed: passed, Weight: weight, Detail: detail}
}

func grade(total, max int) string {
	if max == 0 {
		return "F"
	}
	pct := float64(total) / float64(max) * 100
	switch {
	case pct >= 90:
		return "A"
	case pct >= 80:
		return "B"
	case pct >= 70:
		return "C"
	case pct >= 55:
		return "D"
	default:
		return "F"
	}
}

func summarize(total, max int, checks []model.CheckResult) string {
	var missing []string
	for _, c := range checks {
		if !c.Passed {
			missing = append(missing, c.ID)
		}
	}
	if len(missing) == 0 {
		return "Excellent README coverage."
	}
	return "Improve: " + strings.Join(missing, ", ")
}

// HasContributing is exported for future scorecard extensions.
func HasContributing(md string) bool {
	return reContributing.MatchString(md)
}
