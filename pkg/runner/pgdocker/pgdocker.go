package pgdocker

import (
	"database/sql"
	"fmt"

	"github.com/fahmifan/tblsrun"
	"github.com/fahmifan/tblsrun/pkg/dbtool"
	"github.com/fahmifan/tblsrun/pkg/runner"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var _ runner.DbDriver = (*PostgresDocker)(nil)

const (
	defaultSchema = "public"
	defaultDB     = "postgres"
)

type PostgresDocker struct {
	cfg   tblsrun.Config
	dbCfg tblsrun.Database

	pool     *dockertest.Pool
	resource *dockertest.Resource
}

func NewPostgresDocker(cfg tblsrun.Config) *PostgresDocker {
	return &PostgresDocker{
		cfg: cfg,
	}
}

func (pd *PostgresDocker) DSN() string {
	return pd.dbCfg.
		WithDBName(pd.cfg.TBLS.DBName).
		WithSchema(pd.cfg.TBLS.Schema).
		DSN()
}

func (pd *PostgresDocker) WithSchema(schema string) runner.DbDriver {
	newP := *pd

	newP.cfg.TBLS.Schema = schema
	return &newP
}

func (pd *PostgresDocker) Stop() error {
	return pd.pool.Purge(pd.resource)
}

func (pd *PostgresDocker) Init() (err error) {
	dbcfg := tblsrun.Database{
		Username: "postgres",
		Password: "postgres",
		Host:     "localhost",
		Port:     pd.cfg.TBLS.DBPort,
	}

	pd.pool, err = dockertest.NewPool("")
	if err != nil {
		return fmt.Errorf("init pool: %w", err)
	}

	pd.resource, err = pd.pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "13",
		Env: []string{
			(fmt.Sprintf("POSTGRES_USER=%s", dbcfg.Username)),
			(fmt.Sprintf("POSTGRES_DB=%s", dbcfg.Name)),
			fmt.Sprintf("POSTGRES_PASSWORD=%s", dbcfg.Password),
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
		return fmt.Errorf("run postgres: %w", err)
	}
	dbcfg.Port = pd.resource.GetPort("5432/tcp")

	fmt.Println("try postgres connection")
	var db *sql.DB
	err = pd.pool.Retry(func() (err error) {
		db, err = sql.Open("postgres", dbcfg.DSN())
		if err != nil {
			return fmt.Errorf("open: %w", err)
		}
		defer dbtool.WrapClose(db.Close)
		if err = db.Ping(); err != nil {
			return fmt.Errorf("ping: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("try connection: %w", err)
	}

	pd.dbCfg = dbcfg
	return nil
}

func (pd *PostgresDocker) CreateDB() error {
	if pd.isDefaultDB() {
		return nil
	}

	db, err := dbtool.OpenDB(pd.dbCfg.DSN())
	if err != nil {
		return fmt.Errorf("open default db: %w", err)
	}
	defer dbtool.WrapClose(db.Close)

	// create target database
	if err = dbtool.CreateDB(db, pd.cfg.TBLS.DBName); err != nil {
		return fmt.Errorf("create db: %w", err)
	}

	return nil
}

func (pd *PostgresDocker) CreateSchema() error {
	dsn := pd.dbCfg.DSN()
	if !pd.isDefaultDB() {
		dsn = pd.dbCfg.WithDBName(pd.cfg.TBLS.DBName).DSN()
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}
	defer dbtool.WrapClose(db.Close)

	if pd.cfg.TBLS.Schema != defaultSchema {
		fmt.Println("creating schema")
		if err = dbtool.CreateSchemaIfNoExist(db, pd.cfg.TBLS.Schema); err != nil {
			return fmt.Errorf("create schema: %w", err)
		}
	}

	return nil
}

func (p *PostgresDocker) CreateSchemas() error {
	dsn := p.dbCfg.DSN()
	if !p.isDefaultDB() {
		dsn = p.dbCfg.
			WithDBName(p.cfg.TBLS.DBName).
			DSN()
	}

	db, err := dbtool.OpenDB(dsn)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}

	for _, schema := range p.cfg.TBLS.GetSchemas() {
		if err = dbtool.CreateSchemaIfNoExist(db, schema); err != nil {
			return fmt.Errorf("create schema: %w", err)
		}
		fmt.Println("created schema", schema)
	}

	return nil
}

func (pd *PostgresDocker) isDefaultDB() bool {
	return pd.cfg.TBLS.DBName == defaultDB
}
