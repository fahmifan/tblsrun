package runner

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/fahmifan/tblsrun"
)

type DbDriver interface {
	Init() error
	CreateDB() error
	CreateSchema() error
	DSN() string
	Stop() error
}

type DbMigrator func(dsn, migrationDir string) error

type Runner struct {
	cfg       tblsrun.Config
	dbDriver  DbDriver
	migrateDB DbMigrator
}

func NewRunner(
	cfg tblsrun.Config,
	dbDriver DbDriver,
	migrateDB DbMigrator,
) *Runner {
	return &Runner{
		cfg:       cfg,
		dbDriver:  dbDriver,
		migrateDB: migrateDB,
	}
}

func (r *Runner) Run() error {
	return r.run()
}

func (r *Runner) run() error {
	if _, err := installTblsIfNotExists(); err != nil {
		return err
	}

	if err := r.dbDriver.Init(); err != nil {
		return fmt.Errorf("init: %w", err)
	}
	defer closErr(r.dbDriver.Stop)

	if err := r.dbDriver.CreateDB(); err != nil {
		return fmt.Errorf("create db: %w", err)
	}

	if err := r.dbDriver.CreateSchema(); err != nil {
		return fmt.Errorf("create schema: %w", err)
	}

	if err := r.migrateDB(r.dbDriver.DSN(), r.cfg.TBLS.MigrationDir); err != nil {
		return fmt.Errorf("migrate db: %w", err)
	}

	if err := generateDoc(r.dbDriver.DSN(), r.cfg.TBLS.CfgFile, os.Stdout); err != nil {
		return fmt.Errorf("generate doc: %w", err)
	}

	return nil
}

func installTblsIfNotExists() (path string, err error) {
	if path, err = exec.LookPath("tbls"); err == nil {
		return "", nil
	}

	// try installing it
	fmt.Println("tbls not found. Installing it")

	cmd := exec.Command("go", "install", "github.com/k1LoW/tbls@main")
	out, err := cmd.CombinedOutput()
	fmt.Println(cmd.String())
	fmt.Println(string(out))
	if err != nil {
		return "", err
	}
	return path, nil
}

func generateDoc(dsn string, tblCfgFile string, logger io.Writer) error {
	//nolint:gosec
	cmd := exec.Command("tbls", "doc", dsn, "--force", "--config="+tblCfgFile)
	cmd.Stdout = logger
	cmd.Stderr = logger
	return cmd.Run()
}

func closErr(fn func() error) {
	if fn == nil {
		return
	}
	if err := fn(); err != nil {
		fmt.Fprintf(os.Stderr, "error closing: %v", err)
	}
}
