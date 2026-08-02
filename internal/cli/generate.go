package cli

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/kazhetelabs/devfolio-cli/internal/generate"
	"github.com/spf13/cobra"
)

func newGenerateCmd() *cobra.Command {
	var (
		outDir   string
		maxRepos int
		token    string
		timeout  time.Duration
	)

	cmd := &cobra.Command{
		Use:   "generate <github-user>",
		Short: "Fetch profile, score READMEs, emit static portfolio",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithTimeout(cmd.Context(), timeout)
			defer cancel()

			runner := generate.NewRunner(token)
			res, err := runner.Run(ctx, generate.Options{
				Login:    args[0],
				OutDir:   outDir,
				MaxRepos: maxRepos,
				Token:    token,
			})
			if err != nil {
				return err
			}

			fmt.Fprintf(os.Stdout, "Generated portfolio for @%s\n", res.Portfolio.User.Login)
			fmt.Fprintf(os.Stdout, "  repos scored: %d\n", len(res.Portfolio.Repos))
			fmt.Fprintf(os.Stdout, "  avg README:   %.1f / 100\n", res.Portfolio.AverageScore)
			fmt.Fprintf(os.Stdout, "  output:       %s\n", res.OutDir)
			fmt.Fprintf(os.Stdout, "  open:         %s/index.html\n", res.OutDir)
			fmt.Fprintf(os.Stdout, "  scorecard:    %s/scorecard.md\n", res.OutDir)
			return nil
		},
	}

	cmd.Flags().StringVarP(&outDir, "out", "o", "devfolio-out", "output directory")
	cmd.Flags().IntVar(&maxRepos, "max-repos", 12, "max non-fork repos to include")
	cmd.Flags().StringVar(&token, "token", "", "GitHub token (or set GITHUB_TOKEN)")
	cmd.Flags().DurationVar(&timeout, "timeout", 2*time.Minute, "overall timeout")
	return cmd
}
