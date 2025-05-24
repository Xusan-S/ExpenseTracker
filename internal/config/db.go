package config 

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDB() *pgxpool.Pool {
dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Не удалось подключиться к БД: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("БД не отвечает: %v", err)
	}

	log.Println("Подключение к БД установлено")

	if err := AutoMigrate(pool); err != nil {
		log.Fatalf("Ошибка миграции: %v", err)
	}

	return pool
}

func AutoMigrate(db *pgxpool.Pool) error {
	sql := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		role TEXT NOT NULL CHECK (role IN ('user', 'admin')),
		created_at TIMESTAMP DEFAULT now()
	);

	CREATE TABLE IF NOT EXISTS transactions (
		id SERIAL PRIMARY KEY,
		user_id INT REFERENCES users(id),
		type TEXT CHECK (type IN ('income', 'expense')),
		category TEXT NOT NULL,
		amount NUMERIC NOT NULL,
		date DATE NOT NULL,
		description TEXT,
		receipt_path TEXT,
		created_at TIMESTAMP DEFAULT now()
	);
	`
	_, err := db.Exec(context.Background(), sql)
	return err
}
