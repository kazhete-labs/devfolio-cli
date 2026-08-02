package score

import "testing"

func TestScoreREADME_Excellent(t *testing.T) {
	md := `# Devfolio

![ci](https://img.shields.io/badge/ci-passing-brightgreen)

## Install

` + "```" + `
go install github.com/kazhetelabs/devfolio-cli/cmd/devfolio@latest
` + "```" + `

## Demo

![demo](docs/demo.gif)

## Architecture

CLI → GitHub adapter → scoring domain → static emit.

## License

MIT License
`
	s := ScoreREADME(md)
	if s.Total < 80 {
		t.Fatalf("expected high score, got %d/%d grade=%s missing=%s", s.Total, s.Max, s.Grade, s.Summary)
	}
	if s.Grade != "A" && s.Grade != "B" {
		t.Fatalf("expected A/B, got %s", s.Grade)
	}
}

func TestScoreREADME_Empty(t *testing.T) {
	s := ScoreREADME("")
	if s.Total != 0 {
		t.Fatalf("empty should score 0, got %d", s.Total)
	}
	if s.Grade != "F" {
		t.Fatalf("expected F, got %s", s.Grade)
	}
}

func TestScoreREADME_Partial(t *testing.T) {
	md := "# Tiny\n\nShort readme without install demo or license."
	s := ScoreREADME(md)
	if s.Checks[0].Passed != true {
		t.Fatal("non_empty should pass")
	}
	foundFail := false
	for _, c := range s.Checks {
		if c.ID == "install" && !c.Passed {
			foundFail = true
		}
	}
	if !foundFail {
		t.Fatal("install check should fail")
	}
}
