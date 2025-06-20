package usecases

import (
	"context"
	"errors"

	"github.com/JuanPidarraga/talkus-backend/internal/models"
	"github.com/JuanPidarraga/talkus-backend/internal/repositories"
)

type UserUsecase struct {
	repo *repositories.UserRepository
}

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

// EditUserProfile ejecuta la lógica para editar la foto de perfil de un usuario.
func (u *UserUsecase) EditUserProfile(ctx context.Context, userID string, user models.User) error {
	if userID == "" {
		return errors.New("falta el parámetro 'id'")
	}

	err := u.repo.EditUserProfile(ctx, userID, user)
	if err != nil {
		return errors.New("error editando foto de perfil")
	}
	return nil
}
