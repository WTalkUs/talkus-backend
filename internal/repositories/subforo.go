package repositories

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/JuanPidarraga/talkus-backend/internal/models"
	"google.golang.org/api/iterator"
)

type SubforoRepository struct {
	db *firestore.Client
}

func NewSubforoRepository(db *firestore.Client) *SubforoRepository {
	return &SubforoRepository{db: db}
}

// GetSubforoByID obtiene un subforo por su ID
func (r *SubforoRepository) GetSubforoByID(ctx context.Context, id string) (*models.Subforo, error) {
	doc, err := r.db.Collection("subforos").Doc(id).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("error al obtener el subforo por ID %s: %w", id, err)
	}

	var subforo models.Subforo
	if err := doc.DataTo(&subforo); err != nil {
		return nil, fmt.Errorf("error al decodificar el subforo: %w", err)
	}
	subforo.ForumID = doc.Ref.ID

	return &subforo, nil
}

// GetAll obtiene todos los subforos ordenados por fecha de creación
func (r *SubforoRepository) GetAll(ctx context.Context) ([]*models.Subforo, error) {
	iter := r.db.
		Collection("subforos").
		OrderBy("created_at", firestore.Desc).
		Documents(ctx)
	defer iter.Stop()

	subforos := make([]*models.Subforo, 0)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error al iterar subforos: %w", err)
		}

		var subforo models.Subforo
		if err := doc.DataTo(&subforo); err != nil {
			return nil, fmt.Errorf("error al decodificar subforo: %w", err)
		}
		subforo.ForumID = doc.Ref.ID

		subforos = append(subforos, &subforo)
	}

	return subforos, nil
}

// Create crea un nuevo subforo
func (r *SubforoRepository) Create(ctx context.Context, subforo *models.Subforo) error {
	subforo.CreatedAt = time.Now()
	doc, _, err := r.db.Collection("subforos").Add(ctx, map[string]interface{}{
		"title":       subforo.Title,
		"description": subforo.Description,
		"created_by":  subforo.CreatedBy,
		"category":    subforo.Category,
		"moderators":  subforo.Moderators,
		"is_active":   subforo.IsActive,
		"created_at":  subforo.CreatedAt,
	})
	if err != nil {
		return err
	}
	subforo.ForumID = doc.ID
	return nil
}

// Delete elimina un subforo por ID
func (r *SubforoRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Collection("subforos").Doc(id).Delete(ctx)
	if err != nil {
		return fmt.Errorf("error al eliminar el subforo: %w", err)
	}
	return nil
}

// edita un subforo existente
func (r *SubforoRepository) Edit(ctx context.Context, id string, subforo *models.Subforo) error {
	_, err := r.db.Collection("subforos").Doc(id).Set(ctx, map[string]interface{}{
		"title":       subforo.Title,
		"description": subforo.Description,
		"category":    subforo.Category,
		"moderators":  subforo.Moderators,
		"is_active":   subforo.IsActive,
	}, firestore.MergeAll)
	if err != nil {
		return fmt.Errorf("error al editar el subforo: %w", err)
	}
	return nil
}
