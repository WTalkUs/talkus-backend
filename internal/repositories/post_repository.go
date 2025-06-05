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

func (r *PostRepository) GetPostByID(ctx context.Context, id string) (*models.PostWithAuthor, error) {
	doc, err := r.db.Collection("posts").Doc(id).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("error al obtener el post por ID %s: %w", id, err)
	}

	var p models.Post
	if err := doc.DataTo(&p); err != nil {
		return nil, fmt.Errorf("error al decodificar el post: %w", err)
	}
	p.ID = doc.Ref.ID

	post := &models.PostWithAuthor{
		Post: p,
	}

	if p.AuthorID != "" {
		userDoc, err := r.db.Collection("users").Doc(p.AuthorID).Get(ctx)
		if err != nil {
			return nil, fmt.Errorf("error al obtener el autor con ID %s: %w", p.AuthorID, err)
		}

		var u models.User
		if err := userDoc.DataTo(&u); err != nil {
			return nil, fmt.Errorf("error al decodificar datos del usuario %s: %w", p.AuthorID, err)
		}
		u.UID = userDoc.Ref.ID

		post.Author = &u
	}

	return post, nil
}

func (r *PostRepository) GetAll(ctx context.Context) ([]*models.Post, error) {
	iter := r.db.
		Collection("posts").
		OrderBy("created_at", firestore.Desc).
		Documents(ctx)
	defer iter.Stop()

	posts := make([]*models.Post, 0)
	authorIDs := make(map[string]bool)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error al iterar posts: %w", err)
		}

		var p models.Post

		if err := doc.DataTo(&p); err != nil {
			return nil, fmt.Errorf("error al decodificar post: %w", err)
		}
		p.ID = doc.Ref.ID

		// Registramos los IDs de autores para buscarlos después
		if p.AuthorID != "" {
			authorIDs[p.AuthorID] = true
		}

		posts = append(posts, &p)
	}

	authorInfo := make(map[string]*models.User)
	for authorID := range authorIDs {
		userDoc, err := r.db.Collection("users").Doc(authorID).Get(ctx)
		if err != nil {

			fmt.Printf("Error al buscar el autor %s: %v\n", authorID, err)
			continue
		}

		var user models.User
		if err := userDoc.DataTo(&user); err != nil {
			fmt.Printf("Error al decodificar usuario %s: %v\n", authorID, err)
			continue
		}
		user.UID = userDoc.Ref.ID
		authorInfo[authorID] = &user
	}

	for _, post := range posts {
		if author, exists := authorInfo[post.AuthorID]; exists {
			post.Author = author
		}
	}

	return posts, nil
}

func (r *PostRepository) Create(ctx context.Context, p *models.Post) error {
	p.CreatedAt = time.Now()
	doc, _, err := r.db.Collection("posts").Add(ctx, map[string]interface{}{
		"title":     p.Title,
		"content":   p.Content,
		"author_id": p.AuthorID,
		//"tags":      p.Tags,
		"is_flagged": p.IsFlagged,
		//"forum_id":  p.ForumID,
		"likes":      p.Likes,
		"dislikes":   p.Dislikes,
		"image_url":  p.ImageURL,
		"image_id":   p.ImageID,
		"created_at": p.CreatedAt,
	})
	if err != nil {
		return err
	}
	p.ID = doc.ID
	return nil
}

func (r *PostRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Collection("posts").Doc(id).Delete(ctx)
	if err != nil {
		return fmt.Errorf("error al eliminar el post: %w", err)
	}
	return nil
}

func (r *PostRepository) Edit(ctx context.Context, id string, p *models.Post) error {
	_, err := r.db.Collection("posts").Doc(id).Set(ctx, map[string]interface{}{
		"title":   p.Title,
		"content": p.Content,
		//"author_id":  p.AuthorID,
		//"tags":      p.Tags,
		//"forum_id":  p.ForumID,
		"image_id":  p.ImageID,
		"image_url": p.ImageURL,
	}, firestore.MergeAll)
	if err != nil {
		return fmt.Errorf("error al editar el post: %w", err)
	}
	return nil
}
