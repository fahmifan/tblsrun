package dockertbls

import (
	"fmt"

	"github.com/joeshaw/envdecode"
	"github.com/joho/godotenv"
)

type Database struct {
	Name     string `env:"DATABASE_NAME"`
	Schema   string `env:"DATABASE_SCHEMA"`
	Username string
	Password string
	Host     string
	Port     string
}

func (db Database) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s", db.Username, db.Password, db.Host, db.Port, db.Name, db.Schema)
}

type Config struct {
	Database     Database
	MigrationDir string `env:"MIGRATION_DIR"`
	TblsCfgFile  string `env:"TBLS_CONFIG_FILE,default=.tbls.yml"`
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
