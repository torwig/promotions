package mysql

import (
	"context"
	"database/sql"
	"os"
	"time"

	"github.com/torwig/promotions/internal/promotions/repo"
	"github.com/torwig/promotions/internal/promotions/types"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository() (*Repository, error) {
	db, err := newMySQLConnection()
	if err != nil {
		return nil, err
	}

	// settings depend on application and MySQL server
	db.SetConnMaxLifetime(3 * time.Minute)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return &Repository{db: db}, nil
}

func (r *Repository) Get(ctx context.Context, id string) (*types.Promotion, error) {
	var promo types.Promotion

	query := "SELECT * FROM promotions WHERE id=?"

	err := r.db.GetContext(ctx, &promo, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, repo.ErrPromotionNotFound
	} else if err != nil {
		return nil, errors.Wrap(err, "unable to get promotion from db")
	}

	return &promo, nil
}

func newMySQLConnection() (*sqlx.DB, error) {
	config := mysql.Config{
		Net:       "tcp",
		User:      os.Getenv("MYSQL_USER"),
		Passwd:    os.Getenv("MYSQL_PASSWORD"),
		Addr:      os.Getenv("MYSQL_ADDR"),
		DBName:    os.Getenv("MYSQL_DATABASE"),
		ParseTime: true,
	}

	db, err := sqlx.Connect("mysql", config.FormatDSN())
	if err != nil {
		return nil, errors.Wrap(err, "unable to connect to MySQL")
	}

	return db, err
}
