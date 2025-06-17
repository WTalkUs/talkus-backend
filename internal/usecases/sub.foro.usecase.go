package usecases

import (
	"context"
	"fmt"

	"github.com/JuanPidarraga/talkus-backend/internal/models"
	"github.com/JuanPidarraga/talkus-backend/internal/repositories"
)

type SubforoUsecase struct {
	repo *repositories.SubforoRepository
}

func NewSubforoUsecase(repo *repositories.SubforoRepository) *SubforoUsecase {
	return &SubforoUsecase{repo: repo}
}
func (u *SubforoUsecase) GetSubforoByID(ctx context.Context, id string) (*models.Subforo, error) {
	return u.repo.GetSubforoByID(ctx, id)
}

func (u *SubforoUsecase) GetAllSubforos(ctx context.Context) ([]*models.Subforo, error) {
	return u.repo.GetAll(ctx)
}

func (u *SubforoUsecase) CreateSubforo(ctx context.Context, subforo *models.Subforo) (*models.Subforo, error) {
	if subforo.Title == "" || subforo.Description == "" || len(subforo.Categories) == 0 {
		return nil, fmt.Errorf("faltan campos obligatorios para crear el subforo")
	}

	if err := u.repo.Create(ctx, subforo); err != nil {
		return nil, err
	}

	return subforo, nil
}

func (u *SubforoUsecase) DeactivateSubforo(ctx context.Context, id string) error {
	return u.repo.Deactivate(ctx, id)
}

func (u *SubforoUsecase) EditSubforo(ctx context.Context, id string, subforo *models.Subforo) (*models.Subforo, error) {
	return u.repo.EditSubforo(ctx, id, subforo)
}
