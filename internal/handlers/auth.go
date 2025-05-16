package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"online_book_store/internal/models"
)

func (b *BookDB) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		log.Fatalf("Для передачи данных требуется JSON-формат")
	}

	var newUser models.Register

	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		log.Fatalf("Ошибка парсинга JSON: %v", err)
	}

	err = b.repo.Register(newUser)
	if err != nil {
		log.Fatalf("Ошибка регистрации пользователя: %v", err)
	}
	w.WriteHeader(http.StatusCreated)
}

func (b *BookDB) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		log.Printf("Требуется JSON")
	}

	var user models.Login
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("Ошибка парсинга JSON: %v", err)
	}

	err = b.repo.Login(user)
	if err != nil {
		log.Printf("Неправильно введены данные: %v", err)
	}

	w.WriteHeader(http.StatusOK)
}
