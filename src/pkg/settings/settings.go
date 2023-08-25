package settings

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/moura1001/ssl-tracker/src/pkg/util"
)

var (
	DATABASE_HOST     string
	DATABASE_USER     string
	DATABASE_PASSWORD string
	DATABASE_DB_NAME  string
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	DATABASE_HOST = util.GetEnv("DATABASE_HOST", "")
	DATABASE_USER = util.GetEnv("DATABASE_USER", "")
	DATABASE_PASSWORD = util.GetEnv("DATABASE_PASSWORD", "")
	DATABASE_DB_NAME = util.GetEnv("DATABASE_DB_NAME", "")
}
