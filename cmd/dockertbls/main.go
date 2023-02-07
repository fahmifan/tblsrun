// Package main provides tools to generate document from database using tbls
package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/fahmifan/dockertbls"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("err: %v\n", err)
		os.Exit(1)
		return
	}
	fmt.Println("finish")
}

const (
	dbWaitTime = 10 * time.Second
	dbPassword = "test"
)

func run() error {
	cfg, err := dockertbls.NewConfig(".env")
	if err != nil {
		return err
	}

	_, err = installTblsIfNotExists()
	if err != nil {
		return err
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		return err
	}
	pool.MaxWait = dbWaitTime

	if err = pool.Client.Ping(); err != nil {
		return err
	}

	fmt.Println("init DB")
	resource, err := initDB(pool, &cfg)
	if err != nil {
		return err
	}
	fmt.Println("success init DB")

	fmt.Println("run migration")
	if err = migrateDB(cfg.Database, cfg.MigrationDir); err != nil {
		return err
	}
	fmt.Println("finish migration")

	fmt.Println("run tbls")
	out, err := generateDoc(cfg.Database, cfg.TblsCfgFile)
	fmt.Println(out) // print std out & stderr
	if err != nil {
		return err
	}
	fmt.Println("finished run tbls")

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		return err
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

func initDB(pool *dockertest.Pool, cfg *dockertbls.Config) (resource *dockertest.Resource, err error) {
	fmt.Println("run postgres in docker")

	// set db config
	cfg.Database.Host = "localhost"
	cfg.Database.Username = "test"
	cfg.Database.Password = "test"

	resource, err = pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "13",
		Env: []string{
			(fmt.Sprintf("POSTGRES_USER=%s", cfg.Database.Username)),
			(fmt.Sprintf("POSTGRES_DB=%s", cfg.Database.Name)),
			fmt.Sprintf("POSTGRES_PASSWORD=%s", dbPassword),
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		return nil, err
	}

	cfg.Database.Port = resource.GetPort("5432/tcp")

	dbCfg := cfg.Database

	const defaultSchema = "public"
	dbCfg.Schema = defaultSchema

	fmt.Println("try postgres connection")
	var db *sql.DB
	err = pool.Retry(func() (err error) {
		db, err = sql.Open("postgres", dbCfg.DSN())
		if err != nil {
			return err
		}
		return db.Ping()
	})
	if err != nil {
		return nil, err
	}

	if cfg.Database.Schema != defaultSchema {
		fmt.Println("creating schema")
		if err = createSchema(db, cfg.Database.Schema); err != nil {
			return nil, err
		}
	}

	return resource, nil
}

func createSchema(db *sql.DB, schema string) (err error) {
	// create & set the real schema
	if _, err = db.Exec(`CREATE SCHEMA ` + schema + `;`); err != nil {
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

func generateDoc(dbCfg dockertbls.Database, tblCfgFile string) (string, error) {
	//nolint:gosec
	cmd := exec.Command("tbls", "doc", dbCfg.DSN(), "--force", "--config="+tblCfgFile)
	out, err := cmd.CombinedOutput()
	return string(out), err
}
