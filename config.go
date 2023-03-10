package tblsrun

import (
	"fmt"
	"strconv"
	"strings"

	_ "github.com/davecgh/go-spew/spew"
	"github.com/joeshaw/envdecode"
	"github.com/joho/godotenv"
)

type Database struct {
	Name     string
	Schema   string
	Username string
	Password string
	Host     string
	Port     string
}

func (d Database) WithSchema(schema string) Database {
	d.Schema = schema
	return d
}

func (d Database) WithDBName(dbName string) Database {
	d.Name = dbName
	return d
}

type TBLS struct {
	DBName       string `env:"TBLS_DATABASE_NAME"`
	DBPort       string `env:"TBLS_DATABASE_PORT,default=5432"`
	Schema       string `env:"TBLS_DATABASE_SCHEMA"`
	MigrationDir string `env:"TBLS_MIGRATION_DIR"`
	CfgFile      string `env:"TBLS_CONFIG_FILE,default=.tbls.yml"`
}

func (tbls TBLS) GetMigrationDirs() []string {
	return strings.Split(tbls.MigrationDir, ",")
}

func (tbls TBLS) GetSchemas() (schemas []string) {
	for _, schema := range strings.Split(tbls.Schema, ",") {
		schemas = append(schemas, strings.TrimSpace(schema))
	}
	return schemas
}

func (tbls TBLS) GetConfigFiles() []string {
	return strings.Split(tbls.CfgFile, ",")
}

func (db Database) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s", db.Username, db.Password, db.Host, db.Port, db.Name, db.Schema)
}

func (db Database) DSNWithoutSchema() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", db.Username, db.Password, db.Host, db.Port, db.Name)
}

// Deprecated: use DSN instead
func (db Database) DSNDefaultDBName() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/postgres?sslmode=disable", db.Username, db.Password, db.Host, db.Port)
}

func (db Database) QuotedName() string {
	return "`" + db.Name + "`"
}

func (db Database) GetPort() uint32 {
	u, _ := strconv.ParseUint(db.Port, 10, 32)
	return uint32(u)
}

type Config struct {
	Database Database
	TBLS     TBLS
}

// NewConfig creates an instance of Config.
// It needs the path of the env file to be used.
func NewConfig(env string) (Config, error) {
	err := godotenv.Load(env)
	if err != nil {
		return Config{}, err
	}

	var config Config
	if err := envdecode.Decode(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}
