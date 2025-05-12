package migrations

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func RunMigrations(db *sql.DB) error {
	err := goose.SetDialect("postgres")
	if err != nil {
		log.Fatalf("Неправильно выбран диалект: %v", err)
	}

	err = goose.Up(db, "../../internal/migrations")
	if err != nil {
		log.Fatalf("Ошибка выполнения миграций: %v", err)
	}
	return nil
}
