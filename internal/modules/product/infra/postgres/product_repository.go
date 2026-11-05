package postgres

import (
	"database/sql"
)

type ProductRepository struct {
	DB *sql.DB
}