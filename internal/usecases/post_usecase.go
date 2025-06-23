package usecases

import (
	"context"

	"github.com/JuanPidarraga/talkus-backend/internal/models"
	"github.com/JuanPidarraga/talkus-backend/internal/repositories"
)

type PostUsecase struct {
	repo *repositories.PostRepository
}

func NewPostUsecase(repo *repositories.PostRepository) *PostUsecase {
	return &PostUsecase{repo: repo}
}

func (u *PostUsecase) GetPostByID(ctx context.Context, id string) (*models.PostWithAuthor, error) {
	post, err := u.repo.GetPostByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (u *PostUsecase) GetAllPosts(ctx context.Context) ([]*models.Post, error) {
	return u.repo.GetAll(ctx)
}

func (u *PostUsecase) CreatePost(ctx context.Context, p *models.Post) (*models.Post, error) {
	if err := u.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (u *PostUsecase) DeletePost(ctx context.Context, id string) error {
	return u.repo.Delete(ctx, id)
}

func (u *PostUsecase) GetPostsByAuthorID(ctx context.Context, authorID string) ([]*models.Post, error) {
	return u.repo.GetPostsByAuthorID(ctx, authorID)
}

func (u *PostUsecase) EditPost(ctx context.Context, id string, p *models.Post) error {
	return u.repo.Edit(ctx, id, p)
}

func (u *PostUsecase) GetPostsILiked(ctx context.Context, userID string) ([]*models.Post, error) {
	return u.repo.GetPostsILiked(ctx, userID)
}

func (u *PostUsecase) SavePost(ctx context.Context, userID, postID string) error {
    return u.repo.SavePostForUser(ctx, userID, postID)
}

func (u *PostUsecase) UnsavePost(ctx context.Context, userID, postID string) error {
    return u.repo.RemoveSavedPost(ctx, userID, postID)
}

func (u *PostUsecase) GetSavedPosts(ctx context.Context, userID string) ([]*models.Post, error) {
    return u.repo.GetSavedPostsByUser(ctx, userID)
}

func (u *PostUsecase) IsPostSaved(ctx context.Context, userID, postID string) (bool, error) {
    return u.repo.IsPostSavedByUser(ctx, userID, postID)
}
