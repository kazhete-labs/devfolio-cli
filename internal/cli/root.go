package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Execute runs the root command.
func Execute() {
	root := &cobra.Command{
		Use:   "devfolio",
		Short: "Turn a GitHub username into a portfolio site + README scorecard",
		Long: `devfolio-cli fetches a public GitHub profile, scores repository READMEs,
and emits a static portfolio folder (HTML + Markdown scorecard).

Set GITHUB_TOKEN for higher API rate limits (optional for public data).`,
	}
	root.AddCommand(newGenerateCmd())
	root.AddCommand(newVersionCmd())

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var version = "0.1.0"

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version)
		},
	}
}
