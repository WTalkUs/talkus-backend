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

func (u *SubforoUsecase) DeactivateSubforo(ctx context.Context, id string) (*models.Subforo, error) {
	err := u.repo.Deactivate(ctx, id)
	if err != nil {
		return nil, err
	}
	return u.repo.GetSubforoByID(ctx, id)
}

func (u *SubforoUsecase) EditSubforo(ctx context.Context, id string, subforo *models.Subforo) (*models.Subforo, error) {
	return u.repo.EditSubforo(ctx, id, subforo)
}

func (u *SubforoUsecase) JoinSubforo(ctx context.Context, subforoID, userID string) error {
	return u.repo.JoinSubforo(ctx, subforoID, userID)
}

func (u *SubforoUsecase) LeaveSubforo(ctx context.Context, subforoID, userID string) error {
	return u.repo.LeaveSubforo(ctx, subforoID, userID)
}

func (u *SubforoUsecase) GetSubforosByUserID(ctx context.Context, userID string) ([]*models.Subforo, error) {
	return u.repo.GetSubforosByUserID(ctx, userID)
}
