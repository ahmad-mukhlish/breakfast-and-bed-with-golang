package driver

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DB struct {
	dbPool *sql.DB
}

var dbConn = &DB{}

const maxOpenDbConn = 10
const maxIdleDbConn = 5

const maxDbLifeTime = 5 * time.Minute

func ConnectSQL(dbPath string) (*DB, error) {
	dbPool, err := NewDatabase(dbPath)
	if err != nil {
		panic(err)
	}

	dbPool.SetMaxOpenConns(maxOpenDbConn)
	dbPool.SetMaxIdleConns(maxIdleDbConn)
	dbPool.SetConnMaxLifetime(maxDbLifeTime)

	dbConn.dbPool = dbPool
	err = testDB(dbConn.dbPool)
	if err != nil {
		return nil, err
	}

	return dbConn, nil

}

func testDB(dbPool *sql.DB) error {
	if err := dbPool.Ping(); err != nil {
		return err
	}
	return nil
}

func NewDatabase(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dbPath)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
