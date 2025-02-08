package store

import (
	"database/sql"
	"errors"
	"time"
)

// Storage defines the methods for storing and retrieving stake data
type Storage interface {
	CreateStake(walletAddress string, amount float64) error
	GetTotalStaked(walletAddress string) (float64, error)
}

// sqliteStorage implements the Storage interface for SQLite database
type sqliteStorage struct {
	db *sql.DB
}

// NewStorage creates a new Storage instance
func NewStorage(db *sql.DB) Storage {
	return &sqliteStorage{db: db}
}

func (s *sqliteStorage) CreateStake(walletAddress string, amount float64) error {
	if walletAddress == "" || amount < 0 {
		return errors.New("invalid parameters")
	}
	query := `INSERT INTO stakes (wallet_address, amount, created_at) VALUES (?, ?, ?)`
	_, err := s.db.Exec(query, walletAddress, amount, time.Now())
	return err
}

func (s *sqliteStorage) GetTotalStaked(walletAddress string) (float64, error) {
	var total float64
	query := `SELECT IFNULL(SUM(amount), 0) FROM stakes WHERE wallet_address = ?`
	err := s.db.QueryRow(query, walletAddress).Scan(&total)
	return total, err
}
