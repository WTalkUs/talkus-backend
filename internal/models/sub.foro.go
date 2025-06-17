package models

import (
	"errors"
	"time"
)

type Subforo struct {
	ForumID     string    `firestore:"-" json:"forumId"`
	Title       string    `firestore:"title" json:"title" validate:"required,min=3,max=100"`
	Description string    `firestore:"description" json:"description" validate:"required,min=10,max=500"`
	CreatedBy   string    `firestore:"created_by" json:"createdBy"`
	CreatedAt   time.Time `firestore:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `firestore:"updated_at" json:"updatedAt"`
	Categories  []string  `firestore:"categories" json:"categories" validate:"required,min=1,max=3"`
	Moderators  []string  `firestore:"moderators" json:"moderators" validate:"required,min=1"`
	IsActive    bool      `firestore:"is_active" json:"isActive"`
}

func (s *Subforo) Validate() error {
	// Validación de título
	if len(s.Title) < 3 || len(s.Title) > 100 {
		return errors.New("el título debe tener entre 3 y 100 caracteres")
	}

	// Validación de categorías
	if len(s.Categories) == 0 {
		return errors.New("debe haber al menos una categoría")
	}
	if len(s.Categories) > 3 {
		return errors.New("no se permiten más de 3 categorías")
	}

	// Validación de descripción
	if len(s.Description) < 10 {
		return errors.New("la descripción debe tener al menos 10 caracteres")
	}

	return nil
}
