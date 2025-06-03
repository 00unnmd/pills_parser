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

func tableExists(db *sql.DB, tableName string) (bool, error) {
	var exists bool
	query := `
			SELECT EXISTS (
            SELECT FROM information_schema.tables 
            WHERE table_schema = 'public' 
            AND table_name = $1
        	)`

	err := db.QueryRow(query, tableName).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("ошибка проверки таблицы: %w", err)
	}
	return exists, nil
}

func createTable(db *sql.DB, name string) error {
	query := fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS %s (
		    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		    pharmacy TEXT,
			region TEXT,
			name TEXT,
			mnn TEXT,
			price FLOAT8 NOT NULL,
			discount FLOAT8 NOT NULL,
			discountPercent FLOAT8 NOT NULL,
			producer TEXT,
			rating FLOAT8 NOT NULL,
			reviewsCount INT NOT NULL,
			searchValue TEXT,
			createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			error TEXT
		)`, name)

	_, err := db.Exec(query)
	return err
}

func saveData(db *sql.DB, tableName string, data []domain.ParsedItem) error {
	exists, err := tableExists(db, tableName)
	if err != nil {
		return fmt.Errorf("ошибка проверки существования таблицы: %w", err)
	}

	if !exists {
		if err := createTable(db, tableName); err != nil {
			return fmt.Errorf("ошибка создания таблицы: %w", err)
		}
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("saveData ошибка старта транзакции: %w", err)
	}
	defer tx.Rollback()

	query := fmt.Sprintf(
		`INSERT INTO %s 
			(pharmacy, region, name, mnn, price, discount, discountPercent, producer, rating, reviewsCount, searchValue, createdAt, error)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`,
		tableName,
	)

	stmt, err := tx.Prepare(query)
	if err != nil {
		return fmt.Errorf("saveData ошибка создания оператора: %w", err)
	}
	defer stmt.Close()

	createdAt := time.Now().Format(time.RFC3339)
	for _, item := range data {
		_, err = stmt.Exec(
			item.Pharmacy,
			item.Region,
			item.Name,
			item.Mnn,
			item.Price,
			item.Discount,
			item.DiscountPercent,
			item.Producer,
			item.Rating,
			item.ReviewsCount,
			item.SearchValue,
			createdAt,
			item.Error,
		)
		if err != nil {
			return fmt.Errorf("saveData ошибка выполнения оператора: %w", err)
		}
	}

	return tx.Commit()
}

func SaveToDB(data []domain.ParsedItem, tableName string) {
	db, err := ConnectToPostgres()
	if err != nil {
		log.Printf("SaveToDB err: %s", err)
		return
	}
	defer db.Close()

	err = saveData(db, tableName, data)
	if err != nil {
		log.Printf("SaveToDB err: %s", err)
		return
	}

	log.Println("Данные успешно сохранены в БД")
}
