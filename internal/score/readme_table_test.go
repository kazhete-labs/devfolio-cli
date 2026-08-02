package score

import "testing"

func TestScoreREADME_Table(t *testing.T) {
	cases := []struct {
		name       string
		md         string
		wantMin    int
		wantMax    int
		mustPass   []string
		mustFail   []string
	}{
		{
			name:     "empty",
			md:       "",
			wantMin:  0,
			wantMax:  0,
			mustFail: []string{"non_empty", "install"},
		},
		{
			name:     "title_only",
			md:       "# Title\n",
			wantMin:  10,
			wantMax:  30,
			mustPass: []string{"non_empty", "title"},
			mustFail: []string{"install", "demo"},
		},
		{
			name: "install_and_license",
			md: `# Tool

## Install

` + "```\n" + `go install example.com/tool@latest
` + "```\n" + `
## License

MIT License
`,
			wantMin:  45,
			wantMax:  80,
			mustPass: []string{"install", "license", "code_samples"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := ScoreREADME(tc.md)
			if s.Total < tc.wantMin || s.Total > tc.wantMax {
				t.Fatalf("score %d not in [%d,%d] summary=%s", s.Total, tc.wantMin, tc.wantMax, s.Summary)
			}
			pass := map[string]bool{}
			for _, c := range s.Checks {
				pass[c.ID] = c.Passed
			}
			for _, id := range tc.mustPass {
				if !pass[id] {
					t.Fatalf("expected %s to pass", id)
				}
			}
			for _, id := range tc.mustFail {
				if pass[id] {
					t.Fatalf("expected %s to fail", id)
				}
			}
		})
	}
}
