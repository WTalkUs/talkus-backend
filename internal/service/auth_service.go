package service

import (
	"context"
	"errors"
	"fmt"

	"firebase.google.com/go/v4/auth"
	"github.com/JuanPidarraga/talkus-backend/config"
	"github.com/JuanPidarraga/talkus-backend/internal/models"
)

type AuthService struct {
	firebase *config.FirebaseApp
}

func NewAuthService(firebase *config.FirebaseApp) *AuthService {
	return &AuthService{
		firebase: firebase,
	}
}
func (s *AuthService) VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	if idToken == "" {
		return nil, errors.New("idToken is empty")
	}

	decodedToken, err := s.firebase.Auth.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, err
	}

	if decodedToken == nil {
		return nil, errors.New("decodedToken is nil")
	}

	fmt.Printf("🔐 Token verificado: %v\n", decodedToken.UID)
	return decodedToken, nil
}

func(s *AuthService) GetUserProfile(ctx context.Context, uid string) (*models.User, error) {
	user, err := s.firebase.Auth.GetUser(ctx, uid)
	if err != nil {
		return nil, err
	}

	return &models.User{
		UID:    user.UID,	
		Email: user.Email,
	}, nil
}