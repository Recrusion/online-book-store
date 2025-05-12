package repositories

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func InitDB() (*sql.DB, error) {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Ошибка загрузки конфига: %v", err)
	}
	connectString := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", connectString)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("База данных не отвечает: %v", err)
	}

	return db, nil
}
