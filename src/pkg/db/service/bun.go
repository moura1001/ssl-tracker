package db_service

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/moura1001/ssl-tracker/src/pkg/util"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var Bun *bun.DB

func Init() {

	DATABASE_HOST := util.GetEnv("DATABASE_HOST", "")
	DATABASE_USER := util.GetEnv("DATABASE_USER", "")
	DATABASE_PASSWORD := util.GetEnv("DATABASE_PASSWORD", "")
	DATABASE_DB_NAME := util.GetEnv("DATABASE_DB_NAME", "")

	dsn := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable", DATABASE_USER, DATABASE_PASSWORD, DATABASE_HOST, DATABASE_DB_NAME)
	pgdb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	Bun = bun.NewDB(pgdb, pgdialect.New())

	if err := Bun.Ping(); err != nil {
		log.Fatalf("error to connect into database. Details: '%s'", err)
	}
}
