package repositories

import (
	"context"
	"fmt"
	"sort"
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
		"title":      p.Title,
		"content":    p.Content,
		"author_id":  p.AuthorID,
		"tags":       p.Tags,
		"is_flagged": p.IsFlagged,
		"forum_id":   p.ForumID,
		"likes":      p.Likes,
		"dislikes":   p.Dislikes,
		"image_url":  p.ImageURL,
		"image_id":   p.ImageID,
		"verdict":    p.Verdict,
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
		"tags": p.Tags,
		//"forum_id":  p.ForumID,
		"image_id":  p.ImageID,
		"image_url": p.ImageURL,
	}, firestore.MergeAll)
	if err != nil {
		return fmt.Errorf("error al editar el post: %w", err)
	}
	return nil
}

func (r *PostRepository) GetPostsByAuthorID(ctx context.Context, authorID string) ([]*models.Post, error) {
	iter := r.db.
		Collection("posts").
		Where("author_id", "==", authorID).
		Documents(ctx)
	defer iter.Stop()

	posts := make([]*models.Post, 0)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error al iterar posts del autor: %w", err)
		}

		var p models.Post
		if err := doc.DataTo(&p); err != nil {
			return nil, fmt.Errorf("error al decodificar post: %w", err)
		}
		p.ID = doc.Ref.ID

		// Obtener información del autor
		if p.AuthorID != "" {
			userDoc, err := r.db.Collection("users").Doc(p.AuthorID).Get(ctx)
			if err == nil {
				var user models.User
				if err := userDoc.DataTo(&user); err == nil {
					user.UID = userDoc.Ref.ID
					p.Author = &user
				}
			}
		}

		posts = append(posts, &p)
	}

	// Ordenar por fecha de creación (más recientes primero)
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].CreatedAt.After(posts[j].CreatedAt)
	})

	return posts, nil
}

func (r *PostRepository) GetPostsILiked(ctx context.Context, userID string) ([]*models.Post, error) {
	// Primero obtener todos los votos del usuario
	votesIter := r.db.
		Collection("votes").
		Where("user_id", "==", userID).
		Documents(ctx)
	defer votesIter.Stop()

	// Recolectar los IDs de posts votados con like
	postIDs := make(map[string]bool)
	for {
		doc, err := votesIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error al iterar votos del usuario: %w", err)
		}

		var vote struct {
			PostID string `firestore:"post_id"`
			Type   string `firestore:"type"`
		}
		if err := doc.DataTo(&vote); err != nil {
			return nil, fmt.Errorf("error al decodificar voto: %w", err)
		}

		// Solo incluir posts que tengan like y no sean comentarios
		if vote.PostID != "" && vote.Type == "like" {
			postIDs[vote.PostID] = true
		}
	}

	// Si no hay posts votados, retornar lista vacía
	if len(postIDs) == 0 {
		return []*models.Post{}, nil
	}

	// Obtener los posts correspondientes
	posts := make([]*models.Post, 0)
	for postID := range postIDs {
		postDoc, err := r.db.Collection("posts").Doc(postID).Get(ctx)
		if err != nil {
			// Si no se puede obtener el post, continuar con el siguiente
			fmt.Printf("Error obteniendo post %s: %v\n", postID, err)
			continue
		}

		var p models.Post
		if err := postDoc.DataTo(&p); err != nil {
			fmt.Printf("Error decodificando post %s: %v\n", postID, err)
			continue
		}
		p.ID = postDoc.Ref.ID

		// Obtener información del autor
		if p.AuthorID != "" {
			userDoc, err := r.db.Collection("users").Doc(p.AuthorID).Get(ctx)
			if err == nil {
				var user models.User
				if err := userDoc.DataTo(&user); err == nil {
					user.UID = userDoc.Ref.ID
					p.Author = &user
				}
			}
		}

		posts = append(posts, &p)
	}

	// Ordenar por fecha de creación (más recientes primero)
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].CreatedAt.After(posts[j].CreatedAt)
	})

	return posts, nil
}

func (r *PostRepository) IncrementReaction(ctx context.Context, postID string, reactionType string, delta int) error {
	field := "likes"
	if reactionType == "dislike" {
		field = "dislikes"
	}
	_, err := r.db.Collection("posts").Doc(postID).Update(ctx, []firestore.Update{
		{Path: field, Value: firestore.Increment(delta)},
	})
	return err
}
