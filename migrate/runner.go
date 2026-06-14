package migrate

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ErfanMohseni20/migration-package/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Runner struct {
	cfg config.Config
	m   *migrate.Migrate
}

func New(cfg config.Config) (*Runner, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	m, err := migrate.New(cfg.MigrationsSource(), cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("create migrator: %w", err)
	}

	return &Runner{cfg: cfg, m: m}, nil
}

func (r *Runner) Close() error {
	if r == nil || r.m == nil {
		return nil
	}

	sourceErr, dbErr := r.m.Close()
	return errors.Join(sourceErr, dbErr)
}

func (r *Runner) Up() error {
	return r.run(func() error {
		return r.m.Up()
	})
}

func (r *Runner) Down(steps int) error {
	if steps <= 0 {
		return fmt.Errorf("steps must be greater than zero")
	}

	return r.run(func() error {
		return r.m.Steps(-steps)
	})
}

func (r *Runner) Steps(steps int) error {
	if steps == 0 {
		return errors.New("steps cannot be zero")
	}

	return r.run(func() error {
		return r.m.Steps(steps)
	})
}

func (r *Runner) Force(version int) error {
	if version < 0 {
		return fmt.Errorf("version must be zero or greater")
	}

	return r.run(func() error {
		return r.m.Force(version)
	})
}

func (r *Runner) Status() (version uint, dirty bool, err error) {
	version, dirty, err = r.m.Version()
	if errors.Is(err, migrate.ErrNilVersion) {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, fmt.Errorf("read migration status: %w", err)
	}

	return version, dirty, nil
}

func Up(cfg config.Config) error {
	runner, err := New(cfg)
	if err != nil {
		return err
	}
	defer runner.Close()

	return runner.Up()
}

func Down(cfg config.Config, steps int) error {
	runner, err := New(cfg)
	if err != nil {
		return err
	}
	defer runner.Close()

	return runner.Down(steps)
}

func Force(cfg config.Config, version int) error {
	runner, err := New(cfg)
	if err != nil {
		return err
	}
	defer runner.Close()

	return runner.Force(version)
}

func Status(cfg config.Config) (uint, bool, error) {
	runner, err := New(cfg)
	if err != nil {
		return 0, false, err
	}
	defer runner.Close()

	return runner.Status()
}

func (r *Runner) run(action func() error) error {
	err := action()
	if err == nil || errors.Is(err, migrate.ErrNoChange) {
		return nil
	}

	if isDirtyError(err) {
		version, dirty, verErr := r.m.Version()
		if verErr == nil {
			return fmt.Errorf(
				"%w (current version=%d dirty=%t; fix with: migrate -action force -version %d)",
				err, version, dirty, version,
			)
		}
	}

	return err
}

func isDirtyError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "Dirty database")
}
