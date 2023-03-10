package runner

import (
	"errors"
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
	CreateSchemas() error
	DSN() string
	Stop() error
	WithSchema(schema string) DbDriver
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
	defer func() {
		panicMsg := recover()
		if panicMsg != nil {
			fmt.Println(panicMsg)
		}
	}()

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

	if err := r.dbDriver.CreateSchemas(); err != nil {
		return fmt.Errorf("create schemas: %w", err)
	}

	schemas := r.cfg.TBLS.GetSchemas()
	migrationDirs := r.cfg.TBLS.GetMigrationDirs()
	cfgFiles := r.cfg.TBLS.GetConfigFiles()

	if len(schemas) != len(migrationDirs) {
		return errors.New("migration dir and schema length must be equal")
	}

	configPairs := make([]struct {
		Schema       string
		MigrationDir string
		CfgFile      string
	}, len(migrationDirs))

	for i := 0; i < len(migrationDirs); i++ {
		configPairs[i].Schema = schemas[i]
		configPairs[i].MigrationDir = migrationDirs[i]
		configPairs[i].CfgFile = cfgFiles[i]
	}

	for _, pair := range configPairs {
		dsn := r.dbDriver.WithSchema(pair.Schema).DSN()
		if err := r.migrateDB(dsn, pair.MigrationDir); err != nil {
			return fmt.Errorf("migrate db: %w", err)
		}
	}

	for _, cfgPair := range configPairs {
		dsn := r.dbDriver.WithSchema(cfgPair.Schema).DSN()
		if err := generateDoc(dsn, cfgPair.CfgFile, os.Stdout); err != nil {
			return fmt.Errorf("generate doc: %w", err)
		}
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
