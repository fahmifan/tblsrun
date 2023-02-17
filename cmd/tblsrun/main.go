package main

import (
	"fmt"
	"os"

	"github.com/fahmifan/tblsrun"
	"github.com/fahmifan/tblsrun/pkg/migrator/gomigrate"
	"github.com/fahmifan/tblsrun/pkg/runner"
	"github.com/fahmifan/tblsrun/pkg/runner/pgdocker"
	"github.com/fahmifan/tblsrun/pkg/runner/pgembed"
	"github.com/spf13/cobra"
)

func main() {
	if err := run(os.Args); err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
		return
	}
}

var cfg tblsrun.Config

func run(args []string) (err error) {
	cfg, err = tblsrun.NewConfig(".env")
	if err != nil {
		return fmt.Errorf("new config: %w", err)
	}

	cmd := &cobra.Command{
		Use:   "tblsrun",
		Short: "Generate database documentation from migration files",
	}

	cmd.SetArgs(args[1:])
	cmd.AddCommand(cmdPostgres())

	return cmd.Execute()
}

func cmdPostgres() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "postgres",
		Short: "Run tbls with postgres",
	}
	cmd.AddCommand(cmdPostgresDocker(), cmdPostgresEmbedded())
	return cmd
}

func cmdPostgresDocker() *cobra.Command {
	return &cobra.Command{
		Use:   "docker",
		Short: "Run tbls with postgres in docker",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runner.
				NewRunner(
					cfg,
					pgdocker.NewPostgresDocker(cfg),
					gomigrate.MigrateFromFile,
				).
				Run()
		},
	}
}

func cmdPostgresEmbedded() *cobra.Command {
	return &cobra.Command{
		Use:   "embedded",
		Short: "Run tbls with embedded postgres",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runner.
				NewRunner(
					cfg,
					pgembed.NewPostgresEmbedded(cfg),
					gomigrate.MigrateFromFile,
				).
				Run()
		},
	}
}
