package main

import (
	"log"
	"net/http"
	"online_book_store/internal/handlers"
	"online_book_store/internal/migrations"
	"online_book_store/internal/repositories"

	"github.com/go-pkgz/routegroup"

	_ "github.com/lib/pq"
)

func main() {
	db, err := repositories.InitDB()
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()

	log.Println("Подключение к базе данных - успешно!")

	err = migrations.RunMigrations(db)
	if err != nil {
		log.Fatalf("Ошибка выполнения миграций: %v", err)
	}

	repo := repositories.NewBookDB(db)
	bookHandler := handlers.NewBookDB(repo)

	mux := http.NewServeMux()
	api := routegroup.Mount(mux, "/api")
	v1 := api.Mount("/v1")
	books := v1.Mount("/books")
	books.HandleFunc("GET /{tablename}", bookHandler.GetAllThings)
	books.HandleFunc("DELETE /delete/{tablename}/{id}", bookHandler.DeleteThings)
	books.HandleFunc("POST /create/{tablename}", bookHandler.CreateThings)
	books.HandleFunc("PATCH /update/{id}/{tablename}", bookHandler.UpdateThings)
	books.HandleFunc("POST /signup", bookHandler.SignUp)
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalf("Ошибка при прослушивании порта: %v", err)
	} else {
		log.Fatalf("Сервер успешно запущен!")
	}

}
