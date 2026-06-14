package main

import (
	"fmt"
	"os"

	"github.com/ErfanMohseni20/migration-package/cli"
	"github.com/ErfanMohseni20/migration-package/config"
)

func main() {
	cfg, err := config.Load(".")
	if err != nil {
		fmt.Fprintf(os.Stderr, "load config: %v\n", err)
		os.Exit(1)
	}

	os.Exit(cli.Run(os.Args[1:], cfg))
}
