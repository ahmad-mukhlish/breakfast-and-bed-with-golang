package db

import (
	"database/sql"
	"github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/internal/repository"
)

type postgresDBRepository struct {
	DB *sql.DB
}

func NewPostgresDBRepository(dbConnection *sql.DB) repository.DatabaseRepository {

	return &postgresDBRepository{
		DB: dbConnection,
	}
}
