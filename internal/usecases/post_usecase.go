package usecases

import (
	"context"
	"time"

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

func (u *PostUsecase) EditPost(ctx context.Context, id string, p *models.Post) error {
    return u.repo.Edit(ctx, id, p)
}

func (u *PostUsecase) ReactPost(
    ctx context.Context,
    userID string,
    postID string,
    reactionType string,
) (*models.Vote, error) {
    if err := u.repo.AddReaction(ctx, postID, reactionType); err != nil {
        return nil, err
    }
    vote := &models.Vote{
        UserID:    userID,
        PostID:    postID,
        CommentID: "", 
        Type:      models.VoteType(reactionType),
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    if err := u.repo.AddVote(ctx, vote); err != nil {
        return nil, err
    }
    return vote, nil
}
