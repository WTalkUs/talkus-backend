package repositories

import (
	"context"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/JuanPidarraga/talkus-backend/internal/models"
	"google.golang.org/api/iterator"
)

type CommentRepository interface {
	CreateComment(ctx context.Context, comment *models.Comment) error
	GetCommentsByPostID(ctx context.Context, postID string) ([]models.Comment, error)
	DeleteComment(ctx context.Context, commentID string) error
	GetCommentByID(ctx context.Context, commentID string) (*models.Comment, error)
}

type commentRepository struct {
	db *firestore.Client
}

// NewCommentRepository crea una nueva instancia de CommentRepository.
func NewCommentRepository(db *firestore.Client) CommentRepository {
	return &commentRepository{db: db}
}

// CreateComment crea un nuevo comentario en la base de datos.
func (r *commentRepository) CreateComment(ctx context.Context, comment *models.Comment) error {
	// Crear un nuevo documento en la colección "comments" y obtener una referencia a él
	docRef := r.db.Collection("comments").NewDoc() // Firestore genera un nuevo ID automáticamente

	// Asignamos el commentID generado por Firestore al comentario
	comment.CommentID = docRef.ID

	// Establecer la fecha de creación
	comment.CreatedAt = time.Now()

	// Usamos docRef para guardar el comentario con el ID generado automáticamente
	_, err := docRef.Set(ctx, comment)
	return err
}

func (r *commentRepository) GetCommentByID(ctx context.Context, commentID string) (*models.Comment, error) {
	doc, err := r.db.Collection("comments").Doc(commentID).Get(ctx)
	if err != nil {
		return nil, err
	}

	var comment models.Comment
	err = doc.DataTo(&comment)
	if err != nil {
		return nil, err
	}

	return &comment, nil
}

// GetCommentsByPostID obtiene todos los comentarios de un post por su ID.
func (r *commentRepository) GetCommentsByPostID(ctx context.Context, postID string) ([]models.Comment, error) {
	var comments []models.Comment
	iter := r.db.Collection("comments").Where("postId", "==", postID).Documents(ctx)

	for {
		doc, err := iter.Next()

		if err == iterator.Done {
			// Si hemos llegado al final de los resultados, salimos del ciclo
			break
		}
		if err != nil {
			// Cualquier otro error lo devolvemos
			return nil, err
		}

		var comment models.Comment
		err = doc.DataTo(&comment)
		if err != nil {
			return nil, err
		}

		comments = append(comments, comment)
	}
	return comments, nil
}

// DeleteComment elimina un comentario por su ID.
func (r *commentRepository) DeleteComment(ctx context.Context, commentID string) error {
	// Asegurarnos de que el commentID no tiene una barra inclinada al final
	commentID = strings.TrimSuffix(commentID, "/") // Esto elimina la barra inclinada al final, si la tiene.

	_, err := r.db.Collection("comments").Doc(commentID).Delete(ctx)
	return err
}
