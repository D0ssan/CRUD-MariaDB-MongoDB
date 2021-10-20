package mariadb_test

import (
	"testing"

	"github.com/d0ssan/CRUD-MariaDB-MongoDB/configs"
	"github.com/d0ssan/CRUD-MariaDB-MongoDB/databases/mariadb"

	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/stretchr/testify/assert"
)

var ConfDB = configs.MariaDB{ // nolint:gochecknoglobals
	Driver:        "mysql",
	Username:      "root",
	Name:          "test_users",
	Host:          "localhost",
	Port:          "3306",
	Password:      "secret",
	PathToMigrate: "migration",
}

func TestConnect(t *testing.T) {
	wrongDNS := configs.MariaDB{
		Driver:        "mysql",
		Username:      "root",
		Name:          "test_users",
		Host:          "localhost",
		Password:      "secret",
		PathToMigrate: "migration",
	}

	wrongDriver := configs.MariaDB{
		Driver:        "postgres",
		Username:      "root",
		Name:          "test_users",
		Host:          "localhost",
		Port:          "3306",
		Password:      "secret",
		PathToMigrate: "databases/mariadb/migrations",
	}

	wrongPathToMigrate := configs.MariaDB{
		Driver:        "mysql",
		Username:      "root",
		Name:          "test_users",
		Host:          "localhost",
		Port:          "3306",
		Password:      "secret",
		PathToMigrate: "WRONG_PATH",
	}

	wrongMigrateUp := configs.MariaDB{
		Driver:        "mysql",
		Username:      "root",
		Name:          "test_users",
		Host:          "localhost",
		Port:          "3306",
		Password:      "secret",
		PathToMigrate: "",
	}

	tt := []struct {
		name string
		cfg  configs.MariaDB
		err  string
	}{
		{"Success connection", ConfDB, ""},
		{"Failed connection: wrong dns", wrongDNS, "could not connect to the database"},
		{"Failed connection: wrong driver", wrongDriver, "could not create migration migration"},
		{"Failed connection: wrong path to migrate", wrongPathToMigrate, "error creating db migration"},
		{"Failed connection: cannot migrate up and down", wrongMigrateUp, "cannot manipulate migration"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			db, err := mariadb.Connect(tc.cfg)
			if tc.err != "" {
				if tc.err == "cannot manipulate migration" {
					assert.NotNil(t, db)
					assert.Error(t, db.Up())
					assert.Error(t, db.Down())
					return
				}
				assert.Nil(t, db)
				assert.Error(t, err, tc.err)
				return
			}
			assert.NotNil(t, db)
			assert.NoError(t, err)
			assert.NoError(t, db.Up())
			assert.NoError(t, db.Down())
		})
	}
}
