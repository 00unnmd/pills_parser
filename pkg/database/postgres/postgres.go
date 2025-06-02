package postgres

import (
	"database/sql"
	"fmt"
	"github.com/00unnmd/pills_parser/internal/domain"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// TODO переписать логику для сохранения в уже существующей бд os.Getenv("DB_NAME")
// каждый раз обновлять таблицы вместо создания новой бд
func loadDBConfig() string {
	return fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable client_encoding=UTF8",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
}

func ConnectToPostgres() (*sql.DB, error) {
	connStr := loadDBConfig()

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к postgres: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("ошибка проверки подключения к postgres: %w", err)
	}

	return db, nil
}

func createTable(db *sql.DB, name string) error {
	query := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			region TEXT,
			name TEXT,
			price FLOAT8 NOT NULL,
			discount FLOAT8 NOT NULL,
			priceOld FLOAT8 NOT NULL,
			maxQuantity INT NOT NULL,
			producer TEXT,
			rating FLOAT8 NOT NULL,
			reviewsCount INT NOT NULL,
			error TEXT
		)
	`, name)

	_, err := db.Exec(query)
	return err
}

func createDBWithTable(db *sql.DB, data map[string][]domain.ParsedItem) (*sql.DB, error) {
	now := time.Now().Format("02_01_2006_1504")
	dbName := "parsing_" + now

	if _, err := db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName)); err != nil {
		return nil, fmt.Errorf("ошибка создания БД. %w", err)
	}

	connStr := loadDBConfig()
	newConnStr := fmt.Sprintf("%s dbname=%s", connStr, dbName)
	newDB, err := sql.Open("postgres", newConnStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к БД. %w", err)
	}

	if err := newDB.Ping(); err != nil {
		newDB.Close()
		return nil, fmt.Errorf("ошибка проверки подключения к БД. %w", err)
	}

	for phName, _ := range data {
		if err := createTable(newDB, phName); err != nil {
			newDB.Close()
			return nil, fmt.Errorf("ошибка создания таблицы: %w", err)
		}
	}

	return newDB, nil
}

func saveData(db *sql.DB, data map[string][]domain.ParsedItem) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("saveData ошибка старта транзакции: %w", err)
	}
	defer tx.Rollback()

	for phKey, phData := range data {
		query := fmt.Sprintf(`
			INSERT INTO %s (region, name, price, discount, priceOld, maxQuantity, producer, rating, reviewsCount, error)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		`, phKey)
		stmt, err := tx.Prepare(query)
		if err != nil {
			return fmt.Errorf("saveData ошибка создания оператора: %w", err)
		}
		defer stmt.Close()

		for _, item := range phData {
			_, err = stmt.Exec(item.Region, item.Name, item.Price, item.Discount, item.PriceOld, item.MaxQuantity, item.Producer, item.Rating, item.ReviewsCount, item.Error)
			if err != nil {
				return fmt.Errorf("saveData ошибка выполнения оператора: %w", err)
			}
		}
	}

	return tx.Commit()
}

func SaveToDB(data map[string][]domain.ParsedItem) {
	db, err := ConnectToPostgres()
	if err != nil {
		log.Printf("InitDB err: %s", err)
		return
	}
	defer db.Close()

	newDB, err := createDBWithTable(db, data)
	if err != nil {
		log.Printf("InitDB err: %s", err)
		return
	}
	defer newDB.Close()

	err = saveData(newDB, data)
	if err != nil {
		log.Printf("InitDB err: %s", err)
		return
	}

	log.Println("Данные успешно сохранены в БД")
}
