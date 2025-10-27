package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

// InitDB inicializa la base de datos PostgreSQL
func InitDB(connectionString string) (*sql.DB, error) {
	// Abrir conexión con PostgreSQL
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	// Verificar que la conexión funcione
	if err = db.Ping(); err != nil {
		return nil, err
	}

	log.Printf("Base de datos PostgreSQL inicializada correctamente")
	return db, nil
}
