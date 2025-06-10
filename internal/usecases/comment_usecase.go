package usecases

import (
	"context"

	"github.com/JuanPidarraga/talkus-backend/internal/models"
	"github.com/JuanPidarraga/talkus-backend/internal/repositories"
)

type CommentUsecase interface {
	CreateComment(ctx context.Context, comment *models.Comment) error
	GetCommentsByPostID(ctx context.Context, postID string) ([]models.Comment, error)
	DeleteComment(ctx context.Context, commentID string) error
	GetCommentByID(ctx context.Context, commentID string) (*models.Comment, error)
}

type commentUsecase struct {
	repo repositories.CommentRepository
}

// NewCommentUsecase crea una nueva instancia de CommentUsecase.
func NewCommentUsecase(repo repositories.CommentRepository) CommentUsecase {
	return &commentUsecase{repo: repo}
}

// CreateComment llama al repositorio para crear un nuevo comentario.
func (uc *commentUsecase) CreateComment(ctx context.Context, comment *models.Comment) error {
	return uc.repo.CreateComment(ctx, comment)
}

// GetCommentByID obtiene un comentario por su ID.
func (u *commentUsecase) GetCommentByID(ctx context.Context, commentID string) (*models.Comment, error) {
	return u.repo.GetCommentByID(ctx, commentID)
}

// GetCommentsByPostID llama al repositorio para obtener los comentarios de un post.
func (uc *commentUsecase) GetCommentsByPostID(ctx context.Context, postID string) ([]models.Comment, error) {
	return uc.repo.GetCommentsByPostID(ctx, postID)
}

// DeleteComment llama al repositorio para eliminar un comentario.
func (uc *commentUsecase) DeleteComment(ctx context.Context, commentID string) error {
	return uc.repo.DeleteComment(ctx, commentID)
}
