package main

import (
	"fmt"
	"time"

	"github.com/kazhetelabs/devfolio-cli/internal/score"
)

func main() {
	md := "# Title\n\n## Install\n\n" + "```\ngo install x\n```\n\n## Demo\n\n![d](d.gif)\n\n## License\n\nMIT License\n\n## Architecture\n\nflow\n![ci](https://img.shields.io/badge/x-y-blue)\n"
	for len(md) < 500 {
		md += " word"
	}
	n := 20000
	start := time.Now()
	for i := 0; i < n; i++ {
		_ = score.ScoreREADME(md)
	}
	fmt.Printf("%.6f\n", time.Since(start).Seconds())
}
