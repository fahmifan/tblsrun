package main

import (
	"fmt"
	"os"

	"github.com/fahmifan/tblsrun"
	"github.com/fahmifan/tblsrun/pkg/migrator/gomigrate"
	"github.com/fahmifan/tblsrun/pkg/runner"
	"github.com/fahmifan/tblsrun/pkg/runner/pgembed"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
		return
	}
}

func run() error {
	cfg, err := tblsrun.NewConfig(".env")
	if err != nil {
		return fmt.Errorf("new config: %w", err)
	}

	runs := runner.NewRunner(
		cfg,
		pgembed.NewPostgresEmbedded(cfg),
		gomigrate.MigrateFromFile,
	)

	return runs.Run()
}
