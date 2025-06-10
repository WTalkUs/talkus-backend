package repositories

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/JuanPidarraga/talkus-backend/internal/models"
	"google.golang.org/api/iterator"
)

type VoteRepository interface {
	CreateVote(ctx context.Context, vote *models.Vote) error
	GetVoteByID(ctx context.Context, voteID string) (*models.Vote, error)
	GetVotesByPostID(ctx context.Context, postID string) ([]models.Vote, error)
	GetVotesByCommentID(ctx context.Context, commentID string) ([]models.Vote, error)
	DeleteVote(ctx context.Context, voteID string) error
}

type voteRepository struct {
	db *firestore.Client
}

func NewVoteRepository(db *firestore.Client) VoteRepository {
	return &voteRepository{db: db}
}

func (r *voteRepository) CreateVote(ctx context.Context, vote *models.Vote) error {
	// Validar que solo uno de postId o commentId esté presente
	if vote.PostID == "" && vote.CommentID == "" {
		return fmt.Errorf("un voto debe estar asociado a un post o un comentario")
	}

	if vote.VoteID == "" {
		docRef := r.db.Collection("votes").NewDoc() // Firestore genera un nuevo ID automáticamente
		vote.VoteID = docRef.ID                     // Asignamos el ID generado por Firestore al voto
	}

	vote.CreatedAt = time.Now()

	// Crear o actualizar el voto en la colección de Firestore
	_, err := r.db.Collection("votes").Doc(vote.VoteID).Set(ctx, vote)

	return err
}

func (r *voteRepository) GetVoteByID(ctx context.Context, voteID string) (*models.Vote, error) {
	doc, err := r.db.Collection("votes").Doc(voteID).Get(ctx)
	if err != nil {
		return nil, err
	}

	var vote models.Vote
	err = doc.DataTo(&vote)
	if err != nil {
		return nil, err
	}

	return &vote, nil
}

func (r *voteRepository) GetVotesByPostID(ctx context.Context, postID string) ([]models.Vote, error) {
	var votes []models.Vote
	iter := r.db.Collection("votes").Where("postId", "==", postID).Documents(ctx)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var vote models.Vote
		err = doc.DataTo(&vote)
		if err != nil {
			return nil, err
		}

		votes = append(votes, vote)
	}
	return votes, nil
}

func (r *voteRepository) GetVotesByCommentID(ctx context.Context, commentID string) ([]models.Vote, error) {
	var votes []models.Vote
	iter := r.db.Collection("votes").Where("commentId", "==", commentID).Documents(ctx)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var vote models.Vote
		err = doc.DataTo(&vote)
		if err != nil {
			return nil, err
		}

		votes = append(votes, vote)
	}
	return votes, nil
}

func (r *voteRepository) DeleteVote(ctx context.Context, voteID string) error {
	_, err := r.db.Collection("votes").Doc(voteID).Delete(ctx)
	return err
}
