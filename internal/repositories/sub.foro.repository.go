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
		Where("is_active", "==", true).
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
		"categories":  subforo.Categories,
		"updated_at":  subforo.CreatedAt,
		"moderators":  subforo.Moderators,
		"is_active":   subforo.IsActive,
		"created_at":  subforo.CreatedAt,
		"banner_url":  subforo.BannerURL,
		"icon_url":    subforo.IconURL,
		"members":     subforo.Members,
	})
	if err != nil {
		return err
	}
	subforo.ForumID = doc.ID
	return nil
}

func (r *SubforoRepository) Deactivate(ctx context.Context, id string) error {

	_, err := r.db.Collection("subforos").Doc(id).Set(ctx, map[string]interface{}{
		"is_active":  false,
		"updated_at": time.Now(),
	}, firestore.MergeAll)
	if err != nil {
		return fmt.Errorf("error al desactivar el subforo: %w", err)
	}
	return err
}

// edita un subforo existente
func (r *SubforoRepository) EditSubforo(ctx context.Context, id string, subforo *models.Subforo) (*models.Subforo, error) {
	_, err := r.db.Collection("subforos").Doc(id).Set(ctx, map[string]interface{}{
		"title":       subforo.Title,
		"description": subforo.Description,
		"categories":  subforo.Categories,
		"banner_url":  subforo.BannerURL,
		"icon_url":    subforo.IconURL,
		"moderators":  subforo.Moderators,
		"is_active":   subforo.IsActive,
		"updated_at":  time.Now(),
	}, firestore.MergeAll)

	if err != nil {
		return nil, fmt.Errorf("error al editar el subforo: %w", err)
	}

	updatedSubforo, err := r.GetSubforoByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error al obtener el subforo actualizado: %w", err)
	}

	return updatedSubforo, nil
}

// Unir usuario a subforo
func (r *SubforoRepository) JoinSubforo(ctx context.Context, subforoID, userID string) error {
	_, err := r.db.Collection("subforos").Doc(subforoID).Update(ctx, []firestore.Update{
		{Path: "members", Value: firestore.ArrayUnion(userID)},
		{Path: "updated_at", Value: time.Now()},
	})
	return err
}

// Salir de subforo
func (r *SubforoRepository) LeaveSubforo(ctx context.Context, subforoID, userID string) error {
	_, err := r.db.Collection("subforos").Doc(subforoID).Update(ctx, []firestore.Update{
		{Path: "members", Value: firestore.ArrayRemove(userID)},
		{Path: "updated_at", Value: time.Now()},
	})
	return err
}

func (r *SubforoRepository) GetSubforosByUserID(ctx context.Context, userID string) ([]*models.Subforo, error) {
	// Crear un mapa para evitar duplicados
	subforosMap := make(map[string]*models.Subforo)

	// Consulta 1: Subforos donde el usuario es miembro
	membersIter := r.db.Collection("subforos").Where("members", "array-contains", userID).Documents(ctx)
	defer membersIter.Stop()

	for {
		doc, err := membersIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error al iterar subforos donde es miembro: %w", err)
		}

		var subforo models.Subforo
		if err := doc.DataTo(&subforo); err != nil {
			return nil, fmt.Errorf("error al decodificar subforo: %w", err)
		}
		subforo.ForumID = doc.Ref.ID
		subforosMap[doc.Ref.ID] = &subforo
	}

	// Consulta 2: Subforos donde el usuario es moderador
	moderatorsIter := r.db.Collection("subforos").Where("moderators", "array-contains", userID).Documents(ctx)
	defer moderatorsIter.Stop()

	for {
		doc, err := moderatorsIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error al iterar subforos donde es moderador: %w", err)
		}

		var subforo models.Subforo
		if err := doc.DataTo(&subforo); err != nil {
			return nil, fmt.Errorf("error al decodificar subforo: %w", err)
		}
		subforo.ForumID = doc.Ref.ID
		subforosMap[doc.Ref.ID] = &subforo
	}

	// Convertir el mapa a slice
	subforos := make([]*models.Subforo, 0, len(subforosMap))
	for _, subforo := range subforosMap {
		subforos = append(subforos, subforo)
	}

	return subforos, nil
}
