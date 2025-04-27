package repositories

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/JuanPidarraga/talkus-backend/internal/models"
	"google.golang.org/api/iterator"
)

type PostRepository struct {
	db *firestore.Client
}

func NewPostRepository(db *firestore.Client) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) GetAll(ctx context.Context) ([]*models.Post, error) {
	iter := r.db.
		Collection("posts").
		OrderBy("created_at", firestore.Desc).
		Documents(ctx)
	defer iter.Stop()

	posts := make([]*models.Post, 0)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error iterating posts: %w", err)
		}

		var p models.Post

		if err := doc.DataTo(&p); err != nil {
			return nil, fmt.Errorf("error decoding post: %w", err)
		}
		p.ID = doc.Ref.ID

		posts = append(posts, &p)
	}
	return posts, nil
}

func (r *PostRepository) Create(ctx context.Context, p *models.Post) error {
	p.CreatedAt = time.Now()
	doc, _, err := r.db.Collection("posts").Add(ctx, map[string]interface{}{
		"title":   p.Title,
		"content": p.Content,
		//"author_id": p.AuthorID,
		//"tags":      p.Tags,
		"is_flagged": p.IsFlagged,
		//"forum_id":  p.ForumID,
		"likes":     p.Likes,
		"dislikes":  p.Dislikes,
		"imageUrl":  p.ImageURL,
		"createdAt": p.CreatedAt,
	})
	if err != nil {
		return err
	}
	p.ID = doc.ID
	return nil
}
