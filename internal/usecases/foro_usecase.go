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

// GetAllSubforos obtiene todos los subforos
func (u *SubforoUsecase) GetAllSubforos(ctx context.Context) ([]*models.Subforo, error) {
	return u.repo.GetAll(ctx)
}

// CreateSubforo crea un nuevo subforo
func (u *SubforoUsecase) CreateSubforo(ctx context.Context, subforo *models.Subforo) (*models.Subforo, error) {
	// aqui es agrogo una validacion mas por si acaso
	if subforo.Title == "" || subforo.Description == "" || subforo.Category == "" {
		return nil, fmt.Errorf("faltan campos obligatorios para crear el subforo")
	}

	if err := u.repo.Create(ctx, subforo); err != nil {
		return nil, err
	}

	return subforo, nil
}

// DeleteSubforo elimina un subforo por su ID
// por hora el delete borra el sub for comopletramente
func (u *SubforoUsecase) DeleteSubforo(ctx context.Context, id string) error {
	return u.repo.Delete(ctx, id)
}

// EditSubforo edita un subforo existente
func (u *SubforoUsecase) EditSubforo(ctx context.Context, id string, subforo *models.Subforo) error {
	//aca tambien puedo agrear validaciones adicionales
	return u.repo.Edit(ctx, id, subforo)
}
