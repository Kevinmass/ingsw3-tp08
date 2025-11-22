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

	// Crear las tablas
	if err = createTables(db); err != nil {
		return nil, err
	}

	log.Printf("Base de datos PostgreSQL inicializada correctamente")
	return db, nil
}

// createTables crea el schema de la base de datos PostgreSQL
func createTables(db *sql.DB) error {
	schema := `
	-- Tabla de usuarios
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		username TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	-- Tabla de posts
	CREATE TABLE IF NOT EXISTS posts (
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	-- Tabla de comentarios
	CREATE TABLE IF NOT EXISTS comments (
		id SERIAL PRIMARY KEY,
		post_id INTEGER NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
		user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		content TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	-- Índices para mejorar rendimiento
	CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts(user_id);
	CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments(post_id);
	CREATE INDEX IF NOT EXISTS idx_comments_user_id ON comments(user_id);
	`

	_, err := db.Exec(schema)
	return err
}
