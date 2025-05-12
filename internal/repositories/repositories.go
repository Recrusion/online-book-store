package repositories

import (
	"database/sql"
	"fmt"
	"log"
	"online_book_store/internal/models"
	"strings"
	"sync"
)

type BookDB struct {
	db           *sql.DB
	tableColumns map[string]map[string]bool
	tableOnce    map[string]*sync.Once
	mu           sync.RWMutex
}

func NewBookDB(db *sql.DB) *BookDB {
	b := &BookDB{
		db:           db,
		tableColumns: make(map[string]map[string]bool),
		tableOnce:    make(map[string]*sync.Once),
		mu:           sync.RWMutex{},
	}
	return b
}

func (b *BookDB) DeleteThings(tableName string, id int) (int64, error) {
	exists, err := b.tableExists(tableName)
	if err != nil {
		log.Fatalf("Ошибка проверки наличия таблицы %s: %v", tableName, err)
	}
	if !exists {
		log.Fatalf("Таблицы %s не существует", tableName)
	}

	query := fmt.Sprintf(
		"delete from %s where id = %d",
		tableName,
		id,
	)

	result, err := b.db.Exec(query)
	if err != nil {
		log.Fatalf("Ошибка выполнения SQL-запроса: %v", err)
	}
	return result.RowsAffected()
}

func (b *BookDB) CreateThings(tableName string, data map[string]interface{}) (int64, error) {

	if len(data) == 0 {
		log.Fatalf("Не переданы поля для создания")
	}

	exists, err := b.tableExists(tableName)
	if err != nil {
		log.Fatalf("Ошибка проверки наличия таблицы: %v", err)
	}

	if !exists {
		log.Fatalf("Таблицы с именем %s не существует", tableName)
	}

	for column := range data {
		if !b.isValidColumn(tableName, column) {
			log.Fatalf("Колонки с названием %s в таблице %s не существует", column, tableName)
		}
	}

	setValues := make([]string, 0, len(data))
	args := make([]interface{}, 0, len(data))
	scoreArgs := make([]string, 0, len(data))
	i := 1
	for column, value := range data {
		setValues = append(setValues, column)
		args = append(args, value)
		scoreArgs = append(scoreArgs, fmt.Sprintf("$%d", i))
		i++
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		tableName,
		strings.Join(setValues, ", "),
		strings.Join(scoreArgs, ", "),
	)
	result, err := b.db.Exec(query, args...)
	if err != nil {
		log.Fatalf("Ошибка выполнения SQL-запроса о создании чего-то: %v", err)
	}

	return result.RowsAffected()
}

func (b *BookDB) GetAllThings(tableName string) ([]interface{}, error) {
	exists, err := b.tableExists(tableName)
	if err != nil {
		log.Fatalf("Ошибка проверки наличия таблицы %s: %v", tableName, err)
	}

	if !exists {
		log.Fatalf("Таблицы с названием %s не существует", tableName)
	}
	query := fmt.Sprintf(
		"select * from %s",
		tableName,
	)

	rows, err := b.db.Query(query)
	if err != nil {
		log.Fatalf("Ошибка в получении всех книг: %v", err)
	}
	defer rows.Close()

	var datas []interface{}

	for rows.Next() {
		switch {
		case tableName == "books":
			var data models.Books
			err = rows.Scan(&data.ID, &data.Title, &data.AuthorID, &data.GenreID, &data.Price, &data.StockQuantity)
			if err != nil {
				log.Fatalf("Ошибка в сканировании книг: %v", err)
			}
			datas = append(datas, data)
		case tableName == "genres":
			var data models.Genres
			err = rows.Scan(&data.ID, &data.Name)
			if err != nil {
				log.Fatalf("Ошибка в сканировании книг: %v", err)
			}
			datas = append(datas, data)
		case tableName == "authors":
			var data models.Authors
			err = rows.Scan(&data.ID, &data.Name, &data.Bio)
			if err != nil {
				log.Fatalf("Ошибка в сканировании книг: %v", err)
			}
			datas = append(datas, data)

		}

	}
	return datas, nil
}

func (b *BookDB) UpdateThings(id int, tableName string, updates map[string]interface{}) (int64, error) {
	if len(updates) == 0 {
		log.Fatalf("Не указаны поля для обновления")
	}

	exists, err := b.tableExists(tableName)
	if err != nil {
		log.Fatalf("Ошибка при проверке наличия таблицы: %v", err)
	}

	if !exists {
		log.Fatalf("Таблицы с именем %s не существует", tableName)
	}

	for column := range updates {
		if !b.isValidColumn(tableName, column) {
			log.Fatalf("Колонки с названием %s в таблице %s не существует", column, tableName)
		}
	}

	setValues := make([]string, 0, len(updates))
	args := make([]interface{}, 0, len(updates)+1)
	i := 1
	for column, value := range updates {
		setValues = append(setValues, fmt.Sprintf("%s = $%d", column, i))
		args = append(args, value)
		i++
	}

	args = append(args, id)
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = $%d", tableName, strings.Join(setValues, ", "), i)

	result, err := b.db.Exec(query, args...)
	if err != nil {
		log.Fatalf("Ошибка выполнения SQL запроса: %v", err)
	}

	return result.RowsAffected()
}

func (b *BookDB) isValidColumn(tableName, column string) bool {
	err := b.loadColumns(tableName)
	if err != nil {
		log.Fatalf("Ошибка загрузки колонок для таблицы %s: %v", tableName, err)
	}
	b.mu.Lock()
	result := b.tableColumns[tableName][column]
	defer b.mu.Unlock()
	return result
}

func (b *BookDB) loadColumns(tableName string) error {
	b.mu.Lock()
	_, exists := b.tableOnce[tableName]
	if !exists {
		b.tableOnce[tableName] = &sync.Once{}
	}
	once := b.tableOnce[tableName]
	b.mu.Unlock()
	var err error
	once.Do(func() {
		exists, err := b.tableExists(tableName)
		if err != nil {
			log.Fatalf("Ошибка при проверке наличия таблиц: %v", err)
		}
		if !exists {
			log.Fatalf("Таблица не найдена")
		}

		result, err := b.db.Query(`
			SELECT column_name 
            FROM information_schema.columns 
            WHERE table_name = $1
		`, tableName)
		if err != nil {
			log.Fatalf("Ошибка выполнения запроса о колонках таблицы: %v", err)
		}

		b.mu.Lock()
		b.tableColumns[tableName] = make(map[string]bool)
		b.mu.Unlock()

		for result.Next() {
			var column string
			err = result.Scan(&column)
			if err != nil {
				log.Fatalf("Ошибка получения колонки для таблицы %s: %v", tableName, err)
			}
			b.mu.Lock()
			b.tableColumns[tableName][column] = true
			b.mu.Unlock()
		}

		b.mu.RLock()
		if len(b.tableColumns[tableName]) == 0 {
			log.Fatalf("Не найдены колонки для таблицы %s: %v", tableName, err)
		}
		log.Printf("Загруженные колонки для таблицы %s: %v", tableName, b.tableColumns[tableName])
		b.mu.RUnlock()
	})
	return err
}

func (b *BookDB) tableExists(tableName string) (bool, error) {
	var exists bool
	err := b.db.QueryRow(
		`SELECT EXISTS (
            SELECT FROM information_schema.tables 
            WHERE table_name = $1)`,
		tableName).Scan(&exists)
	if err != nil {
		log.Fatalf("Ошибка выполнения запроса о наличии таблиц в базе данных: %v", err)
	}

	return exists, nil
}
