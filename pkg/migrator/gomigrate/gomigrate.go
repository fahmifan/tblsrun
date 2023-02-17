package gomigrate

import (
	_ "github.com/davecgh/go-spew/spew"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func MigrateFromFile(dsn, srcDir string) error {
	mgr, err := migrate.New("file://"+srcDir, dsn)
	if err != nil {
		return err
	}

	if err = mgr.Up(); err != nil {
		return err
	}

	return nil
}
