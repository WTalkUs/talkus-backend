package repositories

import (
    "context"
    "fmt"

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
        OrderBy("created_at", firestore.Desc). // o "createdAt" si así lo usas
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
        // Rellena Title, Content, AuthorID, CreatedAt
        if err := doc.DataTo(&p); err != nil {
            return nil, fmt.Errorf("error decoding post: %w", err)
        }
        // Asigna el ID que Firestore generó
        p.ID = doc.Ref.ID

        posts = append(posts, &p)
    }
    return posts, nil
}
