package cli

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/ErfanMohseni20/migration-package/config"
	"github.com/ErfanMohseni20/migration-package/migrate"
)

type Options struct {
	Action  string
	Steps   int
	Version int
}

func ParseArgs(args []string) (Options, error) {
	fs := flag.NewFlagSet("migrate", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	action := fs.String("action", "up", "migration action: up, down, force, or status")
	steps := fs.Int("steps", 1, "number of migrations to rollback when action=down")
	version := fs.Int("version", -1, "target version when action=force")

	if err := fs.Parse(args); err != nil {
		return Options{}, err
	}

	return Options{
		Action:  *action,
		Steps:   *steps,
		Version: *version,
	}, nil
}

func Run(args []string, cfg config.Config) int {
	opts, err := ParseArgs(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse flags: %v\n", err)
		return 1
	}

	switch opts.Action {
	case "up":
		err = migrate.Up(cfg)
	case "down":
		err = migrate.Down(cfg, opts.Steps)
	case "force":
		if opts.Version < 0 {
			fmt.Fprintln(os.Stderr, "force requires -version (e.g. -version 2)")
			return 1
		}
		err = migrate.Force(cfg, opts.Version)
	case "status":
		var version uint
		var dirty bool
		version, dirty, err = migrate.Status(cfg)
		if err == nil {
			fmt.Printf("version=%d dirty=%t\n", version, dirty)
			return 0
		}
	default:
		fmt.Fprintf(os.Stderr, "unknown action %q (use up, down, force, or status)\n", opts.Action)
		return 1
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "migration failed: %v\n", err)
		return 1
	}

	if opts.Action != "status" {
		fmt.Println("migration completed")
	}

	return 0
}
