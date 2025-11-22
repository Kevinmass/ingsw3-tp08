package repository

import (
	"database/sql"
	"errors"

	"ingsw3-tp08/internal/models"
)

// PostRepository define las operaciones sobre posts
type PostRepository interface {
	Create(post *models.Post) error
	FindAll() ([]*models.Post, error)
	FindByID(id int) (*models.Post, error)
	Delete(id int) error
	CreateComment(comment *models.Comment) error
	FindCommentsByPostID(postID int) ([]*models.Comment, error)
	DeleteComment(postID int, commentID int, userID int) error
}

// PostgreSQLPostRepository implementa PostRepository usando PostgreSQL
type PostgreSQLPostRepository struct {
	db *sql.DB
}

// NewPostgreSQLPostRepository crea una nueva instancia
func NewPostgreSQLPostRepository(db *sql.DB) *PostgreSQLPostRepository {
	return &PostgreSQLPostRepository{db: db}
}

// Create inserta un nuevo post
func (r *PostgreSQLPostRepository) Create(post *models.Post) error {
	query := `
		INSERT INTO posts (title, content, user_id, created_at)
		VALUES ($1, $2, $3, NOW())
		RETURNING id
	`

	err := r.db.QueryRow(query, post.Title, post.Content, post.UserID).Scan(&post.ID)
	return err
}

// FindAll obtiene todos los posts con información del autor
func (r *PostgreSQLPostRepository) FindAll() ([]*models.Post, error) {
	query := `
		SELECT p.id, p.title, p.content, p.user_id, u.username, p.created_at
		FROM posts p
		JOIN users u ON p.user_id = u.id
		ORDER BY p.created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		post := &models.Post{}
		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.UserID,
			&post.Username,
			&post.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

// FindByID busca un post por ID
func (r *PostgreSQLPostRepository) FindByID(id int) (*models.Post, error) {
	query := `
		SELECT p.id, p.title, p.content, p.user_id, u.username, p.created_at
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.id = $1
	`

	post := &models.Post{}
	err := r.db.QueryRow(query, id).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.UserID,
		&post.Username,
		&post.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return post, nil
}

// Delete elimina un post por ID
func (r *PostgreSQLPostRepository) Delete(id int) error {
	query := `DELETE FROM posts WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// CreateComment inserta un nuevo comentario
func (r *PostgreSQLPostRepository) CreateComment(comment *models.Comment) error {
	query := `
		INSERT INTO comments (post_id, user_id, content, created_at)
		VALUES ($1, $2, $3, NOW())
		RETURNING id
	`

	err := r.db.QueryRow(query, comment.PostID, comment.UserID, comment.Content).Scan(&comment.ID)
	return err
}

// FindCommentsByPostID obtiene todos los comentarios de un post
func (r *PostgreSQLPostRepository) FindCommentsByPostID(postID int) ([]*models.Comment, error) {
	query := `
		SELECT c.id, c.post_id, c.user_id, u.username, c.content, c.created_at
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.post_id = $1
		ORDER BY c.created_at ASC
	`

	rows, err := r.db.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*models.Comment
	for rows.Next() {
		comment := &models.Comment{}
		err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.UserID,
			&comment.Username,
			&comment.Content,
			&comment.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func (r *PostgreSQLPostRepository) DeleteComment(postID int, commentID int, userID int) error {
	query := `
		DELETE FROM comments
		WHERE id = $1 AND post_id = $2 AND user_id = $3
	`
	result, err := r.db.Exec(query, commentID, postID, userID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("no tienes permiso para eliminar este comentario o no existe")
	}
	return nil
}
