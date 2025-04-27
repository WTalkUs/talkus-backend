package repositories

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
)

// UserRepository se encarga de interactuar con la colección "users" en Firestore.
type UserRepository struct {
	db *firestore.Client
}

// NewUserRepository crea una nueva instancia del repositorio.
func NewUserRepository(db *firestore.Client) *UserRepository {
	return &UserRepository{db: db}
}

// GetUserByID busca y retorna el documento del usuario por ID.
func (r *UserRepository) GetUserByID(ctx context.Context, userID string) (map[string]interface{}, error) {
	doc, err := r.db.Collection("users").Doc(userID).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo usuario: %w", err)
	}
	return doc.Data(), nil
}

//CreateUser permite crear un usuario.
func (r *UserRepository) CreateUser(ctx context.Context, userID string, userData map[string]interface{}) error {
	_, err := r.db.Collection("users").Doc(userID).Set(ctx, userData)
	return err
}
