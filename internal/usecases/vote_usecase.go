package usecases

import (
	"context"
	"time"

	"github.com/JuanPidarraga/talkus-backend/internal/models"
	"github.com/JuanPidarraga/talkus-backend/internal/repositories"
)

type VoteUsecase interface {
	CreateVote(ctx context.Context, vote *models.Vote) error
	GetVoteByID(ctx context.Context, voteID string) (*models.Vote, error)
	GetVotesByPostID(ctx context.Context, postID string) ([]models.Vote, error)
	GetVotesByCommentID(ctx context.Context, commentID string) ([]models.Vote, error)
	DeleteVote(ctx context.Context, voteID string) error
	ReactPost(ctx context.Context, userID, postID, reactionType string) (*models.Vote, error)
	GetUserVote(ctx context.Context, userID, postID string) (*models.Vote, error)
}

type voteUsecase struct {
	repo      repositories.VoteRepository
	repoPosts *repositories.PostRepository
}

func NewVoteUsecase(repo repositories.VoteRepository, pr *repositories.PostRepository) VoteUsecase {
	return &voteUsecase{repo: repo,
		repoPosts: pr,}
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

func (u *voteUsecase) ReactPost(
	ctx context.Context,
	userID, postID, reactionType string,
) (*models.Vote, error) {
	// 1) Voto previo
	prev, err := u.repo.GetUserVote(ctx, userID, postID)
	if err != nil {
		return nil, err
	}

	// Anular voto ("none")
	if reactionType == "none" && prev != nil {
		if err := u.repoPosts.IncrementReaction(ctx, postID, string(prev.Type), -1); err != nil {
			return nil, err
		}
		if err := u.repo.DeleteVote(ctx, prev.VoteID); err != nil {
			return nil, err
		}
		return nil, nil
	}

	// 3) Nuevo voto
	if prev == nil {
		if err := u.repoPosts.IncrementReaction(ctx, postID, reactionType, +1); err != nil {
			return nil, err
		}
		v := &models.Vote{
			UserID:    userID,
			PostID:    postID,
			Type:      models.VoteType(reactionType),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := u.repo.CreateVote(ctx, v); err != nil {
			return nil, err
		}
		return v, nil
	}

	// 4) Cambio de voto
	if string(prev.Type) != reactionType {
    // decrementar antiguo
    if err := u.repoPosts.IncrementReaction(ctx, postID, string(prev.Type), -1); err != nil {
        return nil, err
    }
    // incrementar nuevo
    if err := u.repoPosts.IncrementReaction(ctx, postID, reactionType, +1); err != nil {
        return nil, err
    }
    // actualizar vote record
    prev.Type = models.VoteType(reactionType)
    prev.UpdatedAt = time.Now()
    if err := u.repo.DeleteVote(ctx, prev.VoteID); err != nil {
        return nil, err
    }
    if err := u.repo.CreateVote(ctx, prev); err != nil {
        return nil, err
    }
    return prev, nil
}

	// 5) Mismo voto: no hacer nada
	return prev, nil
}

func (u *voteUsecase) GetUserVote(ctx context.Context, userID, postID string) (*models.Vote, error) {
    return u.repo.GetUserVote(ctx, userID, postID)
}


