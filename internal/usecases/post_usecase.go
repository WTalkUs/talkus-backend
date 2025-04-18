package usecases

import(
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

func (u *PostUsecase) GetAllPosts(ctx context.Context) ([]*models.Post, error) {
    return u.repo.GetAll(ctx)
}