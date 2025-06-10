package usecases

import (
	"context"

	"github.com/JuanPidarraga/talkus-backend/internal/models"
	"github.com/JuanPidarraga/talkus-backend/internal/repositories"
)

type VoteUsecase interface {
	CreateVote(ctx context.Context, vote *models.Vote) error
	GetVoteByID(ctx context.Context, voteID string) (*models.Vote, error)
	GetVotesByPostID(ctx context.Context, postID string) ([]models.Vote, error)
	GetVotesByCommentID(ctx context.Context, commentID string) ([]models.Vote, error)
	DeleteVote(ctx context.Context, voteID string) error
}

type voteUsecase struct {
	repo repositories.VoteRepository
}

func NewVoteUsecase(repo repositories.VoteRepository) VoteUsecase {
	return &voteUsecase{repo: repo}
}

func (u *voteUsecase) CreateVote(ctx context.Context, vote *models.Vote) error {
	return u.repo.CreateVote(ctx, vote)
}

func (u *voteUsecase) GetVoteByID(ctx context.Context, voteID string) (*models.Vote, error) {
	return u.repo.GetVoteByID(ctx, voteID)
}

func (u *voteUsecase) GetVotesByPostID(ctx context.Context, postID string) ([]models.Vote, error) {
	return u.repo.GetVotesByPostID(ctx, postID)
}

func (u *voteUsecase) GetVotesByCommentID(ctx context.Context, commentID string) ([]models.Vote, error) {
	return u.repo.GetVotesByCommentID(ctx, commentID)
}

func (u *voteUsecase) DeleteVote(ctx context.Context, voteID string) error {
	return u.repo.DeleteVote(ctx, voteID)
}
