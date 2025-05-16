package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"online_book_store/internal/repositories"
	"strconv"
)

type BookDB struct {
	repo *repositories.BookDB
}

func NewBookDB(repo *repositories.BookDB) *BookDB {
	return &BookDB{repo: repo}
}

func (b *BookDB) GetAllThings(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Fatalf("Неправильный метод запроса")
	}

	tableName := r.PathValue("tablename")

	books, err := b.repo.GetAllThings(tableName)
	if err != nil {
		log.Fatalf("Ошибка выполнение GET запроса: %v", err)
	}
	w.Header().Set("Content-type", "application/json")
	err = json.NewEncoder(w).Encode(books)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Fatalf("Ошибка при переводе в JSON: %v", err)
	}
	w.WriteHeader(http.StatusOK)
}

func (b *BookDB) UpdateThings(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PATCH" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	idCopy := r.PathValue("id")
	if idCopy == "" {
		log.Fatalf("Не передан ID книги")
	}

	tableName := r.PathValue("tablename")
	if tableName == "" {
		log.Fatalf("Не передана таблица для изменения")
	}

	id, err := strconv.Atoi(idCopy)
	if err != nil {
		log.Fatalf("ID не является числом: %v", err)
	}

	if r.Header.Get("Content-Type") != "application/json" {
		log.Fatalf("Для выполнения запроса требуется application/json")
	}

	var updateTable map[string]interface{}

	err = json.NewDecoder(r.Body).Decode(&updateTable)
	if err != nil {
		log.Fatalf("Неправильный формат данных: %v", err)
	}

	_, err = b.repo.UpdateThings(id, tableName, updateTable)
	if err != nil {
		log.Fatalf("Ошибка обновления данных о книге: %v", err)
	}
	w.WriteHeader(http.StatusOK)
}

func (b *BookDB) CreateThings(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Fatalf("Неправильно выбран метод запроса")
	}

	tableName := r.PathValue("tablename")
	if tableName == "" {
		log.Fatalf("Не передана таблица для создания данных")
	}

	if r.Header.Get("Content-Type") != "application/json" {
		log.Fatalf("Для выполнения запроса требуется application/json")
	}

	var data map[string]interface{}

	switch {
	case tableName == "genres":
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			log.Fatalf("Неправильный формат данных")
		}
	case tableName == "books":
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			log.Fatalf("Неправильный формат данных")
		}
	}

	_, err := b.repo.CreateThings(tableName, data)
	if err != nil {
		log.Fatalf("Ошибка при создании новой записи в таблице %s: %v", tableName, err)
	}

	w.WriteHeader(http.StatusOK)
}

func (b *BookDB) DeleteThings(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Fatalf("Неправильно выбран метод запроса")
	}

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Fatalf("Ошибка перевода строки с id в число")
	}

	tableName := r.PathValue("tablename")

	things, err := b.repo.DeleteThings(tableName, id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Fatalf("Ошибка при выполнении DELETE запроса: %v", err)
	}

	err = json.NewEncoder(w).Encode(things)
	if err != nil {
		log.Fatalf("Ошибка перевода в JSON: %v", err)
	}
	w.WriteHeader(http.StatusOK)
}
