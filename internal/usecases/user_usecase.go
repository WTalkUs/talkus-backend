package usecases

import (
	"context"
	"errors"

	"github.com/JuanPidarraga/talkus-backend/internal/repositories"
)

// UserUsecase contiene la lógica de negocio para los usuarios.
type UserUsecase struct {
	repo *repositories.UserRepository
}

// NewUserUsecase crea una nueva instancia de UserUsecase.
func NewUserUsecase(repo *repositories.UserRepository) *UserUsecase {
	return &UserUsecase{repo: repo}
}

// GetUser ejecuta la lógica para obtener un usuario por ID.
func (u *UserUsecase) GetUser(ctx context.Context, userID string) (map[string]interface{}, error) {
	if userID == "" {
		return nil, errors.New("falta el parámetro 'id'")
	}

	user, err := u.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, errors.New("usuario no encontrado")
	}
	return user, nil
}
