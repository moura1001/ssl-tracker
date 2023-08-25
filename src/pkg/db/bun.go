package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/moura1001/ssl-tracker/src/pkg/settings"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var Bun *bun.DB

func Init() {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable", settings.DATABASE_USER, settings.DATABASE_PASSWORD, settings.DATABASE_HOST, settings.DATABASE_DB_NAME)
	pgdb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	Bun = bun.NewDB(pgdb, pgdialect.New())

	if err := Bun.Ping(); err != nil {
		log.Fatalf("error to connect into database. Details: '%s'", err)
	}
}
