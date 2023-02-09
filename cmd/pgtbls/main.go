package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/davecgh/go-spew/spew"
	_ "github.com/davecgh/go-spew/spew"
	"github.com/fahmifan/dockertbls"
	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
		return
	}
}

const (
	defaultSchema = "public"
	defaultDB     = "postgres"
)

func run() error {
	cfg, err := dockertbls.NewConfig(".env")
	if err != nil {
		return err
	}

	cfg.Database.Username = "postgres"
	cfg.Database.Password = "postgres"
	cfg.Database.Host = "localhost"
	cfg.Database.Port = "5432"

	pg := embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().
		Version(embeddedpostgres.V13).
		Port(cfg.Database.GetPort()).
		Database(defaultDB).
		Password(cfg.Database.Password).
		Username(cfg.Database.Username),
	)
	if err = pg.Start(); err != nil {
		return err
	}
	defer logIfErr(pg.Stop)

	time.Sleep(time.Second * 1)
	if err = initDB(cfg); err != nil {
		return err
	}

	spew.Dump(cfg.Database.DSN())
	fmt.Println("run tbls")
	out, err := generateDoc(cfg.Database, cfg.TblsCfgFile)
	fmt.Println(out) // print std out & stderr
	if err != nil {
		return err
	}
	fmt.Println("finished run tbls")

	return nil
}

func generateDoc(dbCfg dockertbls.Database, tblCfgFile string) (string, error) {
	//nolint:gosec
	cmd := exec.Command("tbls", "doc", dbCfg.DSN(), "--force", "--config="+tblCfgFile)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func initDB(cfg dockertbls.Config) error {
	db, err := openDB(cfg.Database.DSNDefaultDBName())
	if err != nil {
		return err
	}
	defer logIfErr(db.Close)

	isDefaultDB := cfg.Database.Name == defaultDB
	if !isDefaultDB {
		fmt.Println("creating database: ", cfg.Database.Name)
		if err = createDB(db, cfg.Database.Name); err != nil {
			return err
		}
	}

	if cfg.Database.Schema != defaultSchema {
		if !isDefaultDB {
			logIfErr(db.Close)

			db, err = openDB(cfg.Database.DSN())
			if err != nil {
				return err
			}
		}

		fmt.Println("creating schema: ", cfg.Database.Schema)
		if err = createSchema(db, cfg.Database.Schema); err != nil {
			return err
		}
	}

	if err = migrateDB(cfg.Database, cfg.MigrationDir); err != nil {
		return err
	}

	return nil
}

func openDB(dsn string) (*sql.DB, error) {
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return gormDB.DB()
}

func createDB(db *sql.DB, dbName string) (err error) {
	if _, err = db.Exec(`CREATE DATABASE "` + dbName + `";`); err != nil {
		return err
	}
	return nil
}

func createSchema(db *sql.DB, schema string) (err error) {
	// create & set the real schema
	if _, err = db.Exec(`CREATE SCHEMA "` + schema + `";`); err != nil {
		return err
	}
	return nil
}

func migrateDB(db dockertbls.Database, migrationDir string) error {
	mgr, err := migrate.New("file://"+migrationDir, db.DSN())
	if err != nil {
		return err
	}

	if err = mgr.Up(); err != nil {
		return err
	}

	return nil
}

func logIfErr(fn func() error) {
	err := fn()
	if err != nil {
		fmt.Println("error: ", err)
	}
}
