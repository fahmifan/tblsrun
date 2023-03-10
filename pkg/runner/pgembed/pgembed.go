package pgembed

import (
	"fmt"
	"io"

	"github.com/fahmifan/tblsrun"
	"github.com/fahmifan/tblsrun/pkg/dbtool"
	"github.com/fahmifan/tblsrun/pkg/runner"
	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var _ runner.DbDriver = (*PostgresEmbedded)(nil)

const (
	defaultSchema = "public"
	defaultDB     = "postgres"
)

type PostgresEmbedded struct {
	pg    *embeddedpostgres.EmbeddedPostgres
	cfg   tblsrun.Config
	dbCfg tblsrun.Database
}

func NewPostgresEmbedded(cfg tblsrun.Config) *PostgresEmbedded {
	return &PostgresEmbedded{
		cfg: cfg,
	}
}

func (p *PostgresEmbedded) DSN() string {
	return p.dbCfg.
		WithDBName(p.cfg.TBLS.DBName).
		WithSchema(p.cfg.TBLS.Schema).
		DSN()
}

func (p *PostgresEmbedded) WithSchema(schema string) runner.DbDriver {
	newP := *p

	newP.cfg.TBLS.Schema = schema
	return &newP
}

func (p *PostgresEmbedded) Stop() error {
	return p.pg.Stop()
}

func (p *PostgresEmbedded) Init() error {
	dbcfg := tblsrun.Database{
		Username: "postgres",
		Password: "postgres",
		Host:     "localhost",
		Port:     p.cfg.TBLS.DBPort,
	}

	p.pg = embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().
		Version(embeddedpostgres.V13).
		Port(dbcfg.GetPort()).
		Database(defaultDB).
		Password(dbcfg.Password).
		Username(dbcfg.Username).
		Logger(io.Discard),
	)
	if err := p.pg.Start(); err != nil {
		return fmt.Errorf("start embedded postgres: %w", err)
	}

	p.dbCfg = dbcfg

	return nil
}

func (p *PostgresEmbedded) CreateDB() error {
	if p.isDefaultDB() {
		return nil
	}

	db, err := dbtool.OpenDB(p.dbCfg.DSN())
	if err != nil {
		return fmt.Errorf("open default db: %w", err)
	}
	defer logIfErr(db.Close)

	// create target database
	if err = dbtool.CreateDB(db, p.cfg.TBLS.DBName); err != nil {
		return fmt.Errorf("create db: %w", err)
	}

	return nil
}

func (p *PostgresEmbedded) CreateSchema() error {
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

	if err = dbtool.CreateSchemaIfNoExist(db, p.cfg.TBLS.Schema); err != nil {
		return fmt.Errorf("create schema: %w", err)
	}

	return nil
}

func (p *PostgresEmbedded) CreateSchemas() error {
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
	}

	return nil
}

func (p *PostgresEmbedded) isDefaultDB() bool {
	return p.cfg.TBLS.DBName == defaultDB
}

func logIfErr(fn func() error) {
	err := fn()
	if err != nil {
		fmt.Println("error: ", err)
	}
}
