package usecases

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/JuanPidarraga/talkus-backend/internal/models"
	"github.com/JuanPidarraga/talkus-backend/internal/repositories"
)

type CommentUsecase interface {
	CreateComment(ctx context.Context, comment *models.Comment) error
	GetCommentsByPostID(ctx context.Context, postID string) ([]models.Comment, error)
	DeleteComment(ctx context.Context, commentID string) error
	GetCommentByID(ctx context.Context, commentID string) (*models.Comment, error)
	UpdateComment(ctx context.Context, commentID string, userID string, updatedContent string) (*models.Comment, error)
	CreateReply(ctx context.Context, parentID string, comment *models.Comment) error
	GetCommentTree(ctx context.Context, postID string) ([]*models.CommentWithReplies, error)
	AddReaction(ctx context.Context, commentID, userID, reaction string) (*models.Comment, error)
}

type commentUsecase struct {
	repo repositories.CommentRepository
}

func NewCommentUsecase(repo repositories.CommentRepository) CommentUsecase {
	return &commentUsecase{repo: repo}
}

func (uc *commentUsecase) CreateComment(ctx context.Context, comment *models.Comment) error {
	if err := comment.Validate(); err != nil {
		return err
	}
	return uc.repo.CreateComment(ctx, comment)
}

func (uc *commentUsecase) GetCommentsByPostID(ctx context.Context, postID string) ([]models.Comment, error) {
	return uc.repo.GetCommentsByPostID(ctx, postID)
}

// GetCommentByID obtiene un comentario por su ID.
func (u *commentUsecase) GetCommentByID(ctx context.Context, commentID string) (*models.Comment, error) {
	return u.repo.GetCommentByID(ctx, commentID)
}

func (uc *commentUsecase) UpdateComment(ctx context.Context, commentID string, userID string, updatedContent string) (*models.Comment, error) {

	if strings.TrimSpace(updatedContent) == "" {
		return nil, fmt.Errorf("comment content cannot be empty")
	}
	if len(updatedContent) > 500 {
		return nil, fmt.Errorf("comment is too long (max 500 chars)")
	}

	comment, err := uc.repo.GetCommentByID(ctx, commentID)
	if err != nil {
		return nil, fmt.Errorf("comment not found")
	}
	if comment.AuthorID != userID {
		return nil, fmt.Errorf("unauthorized: you can only edit your own comments")
	}

	return uc.repo.UpdateComment(ctx, commentID, updatedContent)
}

func (uc *commentUsecase) DeleteComment(ctx context.Context, commentID string) error {
	return uc.repo.DeleteComment(ctx, commentID)
}

func (uc *commentUsecase) CreateReply(ctx context.Context, parentID string, comment *models.Comment) error {
	// 1. Obtener el comentario padre primero
	parent, err := uc.repo.GetCommentByID(ctx, parentID)
	if err != nil {
		return fmt.Errorf("parent comment not found")
	}

	// 2. Asignar el PostID del padre al comentario respuesta
	comment.PostID = parent.PostID
	comment.ParentID = parentID

	// 3. Validar el comentario (ahora incluirá el PostID)
	if err := comment.Validate(); err != nil {
		return err
	}

	// 4. Crear el comentario
	return uc.repo.CreateComment(ctx, comment)
}

func (uc *commentUsecase) GetReplies(ctx context.Context, parentID string) ([]models.Comment, error) {
	// Validación adicional
	if parentID == "" {
		return nil, fmt.Errorf("parentID cannot be empty")
	}

	// Verificar que el comentario padre existe
	if _, err := uc.repo.GetCommentByID(ctx, parentID); err != nil {
		return nil, fmt.Errorf("failed to get parent comment: %w", err)
	}

	return uc.repo.GetReplies(ctx, parentID)
}

func (uc *commentUsecase) GetCommentTree(ctx context.Context, postID string) ([]*models.CommentWithReplies, error) {
	comments, err := uc.repo.GetCommentsByPostID(ctx, postID)
	if err != nil {
		return nil, err
	}

	commentMap := make(map[string]*models.CommentWithReplies)
	var roots []*models.CommentWithReplies

	sort.Slice(comments, func(i, j int) bool {
		return comments[i].CreatedAt.Before(comments[j].CreatedAt)
	})

	for _, c := range comments {
		node := &models.CommentWithReplies{
			Comment: &c,
			Replies: make([]*models.CommentWithReplies, 0),
		}
		commentMap[c.CommentID] = node

		if c.ParentID == "" {
			roots = append(roots, node)
		} else {
			if parent, exists := commentMap[c.ParentID]; exists {
				parent.Replies = append(parent.Replies, node)
			} else {
				// Si el padre no está cargado aún, crear un nodo temporal
				tempParent := &models.CommentWithReplies{
					Comment: &models.Comment{CommentID: c.ParentID},
				}
				commentMap[c.ParentID] = tempParent
				tempParent.Replies = append(tempParent.Replies, node)
			}
		}
	}

	return roots, nil
}

func (uc *commentUsecase) AddReaction(ctx context.Context, commentID, userID, reaction string) (*models.Comment, error) {

	if !(reaction == "like" || reaction == "dislike") {
		return nil, fmt.Errorf("reaction must be either 'like' or 'dislike'")
	}

	if reaction != "like" && reaction != "dislike" {
		return nil, fmt.Errorf("invalid reaction type")
	}

	_, err := uc.repo.GetCommentByID(ctx, commentID)
	if err != nil {
		return nil, fmt.Errorf("comment not found")
	}

	return uc.repo.AddReaction(ctx, commentID, userID, reaction)
}
