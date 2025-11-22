package repository

import (
	"database/sql"

	"ingsw3-tp7-tp8-integrated/backend/internal/models"
)

// UserRepository define las operaciones sobre usuarios
// INTERFACE: permite crear mocks fácilmente para testing
type UserRepository interface {
	Create(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	FindByID(id int) (*models.User, error)
}

// PostgreSQLUserRepository implementa UserRepository usando PostgreSQL
type PostgreSQLUserRepository struct {
	db *sql.DB
}

// NewPostgreSQLUserRepository crea una nueva instancia
func NewPostgreSQLUserRepository(db *sql.DB) *PostgreSQLUserRepository {
	return &PostgreSQLUserRepository{db: db}
}

// Create inserta un nuevo usuario en la base de datos
func (r *PostgreSQLUserRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (email, password, username, created_at)
		VALUES ($1, $2, $3, NOW())
		RETURNING id
	`

	err := r.db.QueryRow(query, user.Email, user.Password, user.Username).Scan(&user.ID)
	return err
}

// FindByEmail busca un usuario por email
func (r *PostgreSQLUserRepository) FindByEmail(email string) (*models.User, error) {
	query := `SELECT id, email, password, username, created_at FROM users WHERE email = $1`

	user := &models.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Username,
		&user.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // Usuario no encontrado (no es error)
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

// FindByID busca un usuario por ID
func (r *PostgreSQLUserRepository) FindByID(id int) (*models.User, error) {
	query := `SELECT id, email, password, username, created_at FROM users WHERE id = $1`

	user := &models.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Username,
		&user.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}
