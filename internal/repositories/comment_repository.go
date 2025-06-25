package repositories

import (
	"context"
	"fmt"
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
	UpdateComment(ctx context.Context, commentID string, updatedContent string) (*models.Comment, error)
	CreateReply(ctx context.Context, parentID string, comment *models.Comment) error
	GetReplies(ctx context.Context, parentID string) ([]models.Comment, error)
	AddReaction(ctx context.Context, commentID, userID, reaction string) (*models.Comment, error)
}

type commentRepository struct {
	db *firestore.Client
}

// NewCommentRepository crea una nueva instancia de CommentRepository.
func NewCommentRepository(db *firestore.Client) CommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) CreateComment(ctx context.Context, comment *models.Comment) error {
	docRef := r.db.Collection("comments").NewDoc()

	comment.CommentID = docRef.ID
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()

	_, err := docRef.Set(ctx, map[string]interface{}{
		"commentID": comment.CommentID,
		"postId":    comment.PostID,
		"authorId":  comment.AuthorID,
		"content":   comment.Content,
		"createdAt": comment.CreatedAt,
		"updatedAt": comment.UpdatedAt,
		"likes":     comment.Likes,
		"dislikes":  comment.Dislikes,
		"parentId":  comment.ParentID,
	})

	if err != nil {
		return err
	}

	userDoc, err := r.db.Collection("users").Doc(comment.AuthorID).Get(ctx)
	if err == nil {
		var user models.User
		if err := userDoc.DataTo(&user); err == nil {
			user.UID = userDoc.Ref.ID
			comment.Author = &user
		}
	}

	return nil
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

func (r *commentRepository) GetCommentsByPostID(ctx context.Context, postID string) ([]models.Comment, error) {

	iter := r.db.Collection("comments").Where("postId", "==", postID).OrderBy("likes", firestore.Desc).Documents(ctx)
	var comments []models.Comment
	var authorIDs []string

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var comment models.Comment
		if err := doc.DataTo(&comment); err != nil {
			continue
		}
		comment.CommentID = doc.Ref.ID
		comments = append(comments, comment)
		authorIDs = append(authorIDs, comment.AuthorID)
	}

	authors, err := r.getUsersByIDs(ctx, authorIDs)
	if err != nil {
		return comments, nil
	}

	for i := range comments {
		if author, exists := authors[comments[i].AuthorID]; exists {
			comments[i].Author = author
		}
	}

	return comments, nil
}

func (r *commentRepository) getUsersByIDs(ctx context.Context, userIDs []string) (map[string]*models.User, error) {
	usersMap := make(map[string]*models.User)

	// Eliminar IDs duplicados
	uniqueIDs := make(map[string]struct{})
	for _, id := range userIDs {
		uniqueIDs[id] = struct{}{}
	}

	// Cargar usuarios desde Firestore
	for id := range uniqueIDs {
		doc, err := r.db.Collection("users").Doc(id).Get(ctx)
		if err != nil {
			continue
		}

		var user models.User
		if err := doc.DataTo(&user); err == nil {
			user.UID = doc.Ref.ID
			usersMap[id] = &user
		}
	}

	return usersMap, nil
}

func (r *commentRepository) UpdateComment(ctx context.Context, commentID string, updatedContent string) (*models.Comment, error) {
	docRef := r.db.Collection("comments").Doc(commentID)

	// 1. Primero obtenemos el comentario existente para preservar el AuthorID
	existingDoc, err := docRef.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("comment not found")
	}

	var existingComment models.Comment
	if err := existingDoc.DataTo(&existingComment); err != nil {
		return nil, fmt.Errorf("failed to parse comment")
	}

	// 2. Actualizamos solo los campos permitidos
	updates := []firestore.Update{
		{Path: "content", Value: updatedContent},
		{Path: "updatedAt", Value: time.Now()},
	}

	if _, err := docRef.Update(ctx, updates); err != nil {
		return nil, fmt.Errorf("failed to update comment")
	}

	// 3. Obtenemos el comentario actualizado
	updatedDoc, err := docRef.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated comment")
	}

	var updatedComment models.Comment
	if err := updatedDoc.DataTo(&updatedComment); err != nil {
		return nil, err
	}
	updatedComment.CommentID = updatedDoc.Ref.ID

	// 4. ¡Cargamos el autor aquí! (igual que en CreateComment)
	userDoc, err := r.db.Collection("users").Doc(existingComment.AuthorID).Get(ctx)
	if err == nil { // Si el usuario existe
		var user models.User
		if err := userDoc.DataTo(&user); err == nil {
			user.UID = userDoc.Ref.ID
			updatedComment.Author = &user
		}
	}

	return &updatedComment, nil
}

func (r *commentRepository) DeleteComment(ctx context.Context, commentID string) error {
	// Asegurarnos de que el commentID no tiene una barra inclinada al final
	commentID = strings.TrimSuffix(commentID, "/") // Esto elimina la barra inclinada al final, si la tiene.

	_, err := r.db.Collection("comments").Doc(commentID).Delete(ctx)
	return err
}

func (r *commentRepository) CreateReply(ctx context.Context, ParentID string, comment *models.Comment) error {
	// Obtener el comentario padre
	parent, err := r.GetCommentByID(ctx, ParentID)
	if err != nil {
		return fmt.Errorf("parent comment not found: %w", err)
	}

	// Asignar los valores necesarios
	comment.PostID = parent.PostID
	comment.ParentID = ParentID // Asegurar que el ParentID está asignado
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()

	// Validar el comentario
	if err := comment.Validate(); err != nil {
		return fmt.Errorf("comment validation failed: %w", err)
	}

	// Crear el documento en Firestore
	docRef := r.db.Collection("comments").NewDoc()
	comment.CommentID = docRef.ID

	// Mapa de datos para Firestore
	data := map[string]interface{}{
		"commentID": comment.CommentID,
		"postId":    comment.PostID,
		"parentId":  comment.ParentID,
		"authorId":  comment.AuthorID,
		"content":   comment.Content,
		"createdAt": comment.CreatedAt,
		"updatedAt": comment.UpdatedAt,
		"likes":     comment.Likes,
		"dislikes":  comment.Dislikes,
	}

	_, err = docRef.Set(ctx, data)
	if err != nil {
		return fmt.Errorf("failed to save comment: %w", err)
	}

	return nil
}

func (r *commentRepository) GetReplies(ctx context.Context, parentID string) ([]models.Comment, error) {
	// 1. Verificar primero que el comentario padre existe
	parentRef := r.db.Collection("comments").Doc(parentID)
	_, err := parentRef.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("parent comment does not exist: %w", err)
	}

	// 2. Realizar la consulta con logs para depuración
	fmt.Printf("Buscando respuestas para parentID: %s\n", parentID)

	iter := r.db.Collection("comments").
		Where("parentId", "==", parentID).
		OrderBy("createdAt", firestore.Asc).
		Documents(ctx)

	defer iter.Stop()

	var replies []models.Comment

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Printf("Error al iterar documentos: %v\n", err)
			return nil, fmt.Errorf("error iterating comments: %w", err)
		}

		fmt.Printf("Documento encontrado ID: %s\n", doc.Ref.ID)

		var comment models.Comment
		if err := doc.DataTo(&comment); err != nil {
			fmt.Printf("Error mapeando documento %s: %v\n", doc.Ref.ID, err)
			continue
		}

		comment.CommentID = doc.Ref.ID
		fmt.Printf("Respuesta encontrada: %+v\n", comment)

		replies = append(replies, comment)
	}

	fmt.Printf("Total respuestas encontradas: %d\n", len(replies))
	return replies, nil
}

func (r *commentRepository) AddReaction(ctx context.Context, commentID, userID, reaction string) (*models.Comment, error) {
	if reaction != "like" && reaction != "dislike" {
		return nil, fmt.Errorf("reacción inválida: %s", reaction)
	}

	docRef := r.db.Collection("comments").Doc(commentID)

	err := r.db.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		doc, err := tx.Get(docRef)
		if err != nil {
			return err
		}

		var comment models.Comment
		if err := doc.DataTo(&comment); err != nil {
			return err
		}

		if comment.Reactions == nil {
			comment.Reactions = make(map[string]string)
		}

		prevReaction, exists := comment.Reactions[userID]

		if !exists {
			// Primera vez que reacciona
			comment.Reactions[userID] = reaction
			if reaction == "like" {
				comment.Likes++
			} else {
				comment.Dislikes++
			}
		} else if prevReaction == reaction {
			// Quitar reacción si es igual a la que ya tiene
			delete(comment.Reactions, userID)
			if reaction == "like" && comment.Likes > 0 {
				comment.Likes--
			}
			if reaction == "dislike" && comment.Dislikes > 0 {
				comment.Dislikes--
			}
		} else {
			// Cambiar reacción
			comment.Reactions[userID] = reaction
			if prevReaction == "like" && comment.Likes > 0 {
				comment.Likes--
			}
			if prevReaction == "dislike" && comment.Dislikes > 0 {
				comment.Dislikes--
			}
			if reaction == "like" {
				comment.Likes++
			}
			if reaction == "dislike" {
				comment.Dislikes++
			}
		}

		return tx.Set(docRef, map[string]interface{}{
			"authorId":  comment.AuthorID,
			"commentID": comment.CommentID,
			"content":   comment.Content,
			"createdAt": comment.CreatedAt,
			"dislikes":  comment.Dislikes,
			"likes":     comment.Likes,
			"reactions": comment.Reactions,
			"postId":    comment.PostID,
			"updatedAt": comment.UpdatedAt,
		})
	})

	if err != nil {
		return nil, err
	}

	return r.GetCommentByID(ctx, commentID)
}
