package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type DBConfigs struct {
	Host     string
	Port     string
	Username string
	DBName   string
	Password string
}

func NewPostgresDB(cfg *DBConfigs) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.Username,
		cfg.Password,
		cfg.DBName,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return db, nil
}
