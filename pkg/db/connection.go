package db

import (
	"context"
	"database/sql"
	"firstapi/configs"
	"fmt"
	"log"
	"time"
)

type Database struct {
	DB *sql.DB
}

func NewConnection(config *configs.DB) (*Database, error) {
	objectconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.NameDB,
		config.SSLMode,
	)

	db, err := sql.Open("postgres", objectconn)
	if err != nil {
		return nil, fmt.Errorf("error opening db connection!!: %w", err)
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	contextdb, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := db.PingContext(contextdb); err != nil {
		db.Close()
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}
	return &Database{DB: db}, nil
}

// Cierra conexion DB
func (conn *Database) Close() error {
	if conn.DB != nil {
		log.Println("Closing database connection...")
		return conn.DB.Close()
	}
	return nil
}

func (conn *Database) Health() error {
	contextdb, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return conn.DB.PingContext(contextdb)
}

func (conn *Database) GetStats() sql.DBStats {
	return conn.DB.Stats()
}

/*MIGRATIONS SE USA CON MIGRATE FROM GO EN LINEA DE COMANDOS*/

func (conn *Database) QueryRow(contextdb context.Context, query string, args ...interface{}) *sql.Row {
	return conn.DB.QueryRowContext(contextdb, query, args...)
}

func (conn *Database) QueryMultipleRows(contextdb context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return conn.DB.QueryContext(contextdb, query, args...)
}

func (conn *Database) ExecSQL(context context.Context, query string, args ...interface{}) (sql.Result, error) {
	return conn.DB.ExecContext(context, query, args...)
}

func (conn *Database) Transaction(context context.Context, fn func(*sql.Tx) error) error {
	tx, err := conn.DB.BeginTx(context, nil)
	if err != nil {
		return fmt.Errorf("error executing database to database: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else {
			err = tx.Commit()
		}
	}()
	err = fn(tx)
	return err
}
