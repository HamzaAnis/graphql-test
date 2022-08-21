package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"time"

	"github.com/hamzaanis/graphql-test/graph/dal"
	"github.com/hamzaanis/graphql-test/graph/errors"
	"github.com/hamzaanis/graphql-test/graph/generated"
	"github.com/hamzaanis/graphql-test/graph/model"
)

// CreateVideo is the resolver for the createVideo field.
func (r *mutationResolver) CreateVideo(ctx context.Context, input model.NewVideo) (*model.Video, error) {
	newVideo := &model.Video{
		URL:       input.URL,
		Name:      input.Name,
		CreatedAt: time.Now().UTC(),
	}

	rows, err := dal.LogAndQuery(r.db, "INSERT INTO videos (name, url, user_id, created_at) VALUES($1, $2, $3, $4) RETURNING id",
		input.Name, input.URL, input.UserID, newVideo.CreatedAt)
	defer rows.Close()

	if err != nil || !rows.Next() {
		return &model.Video{}, err
	}
	if err := rows.Scan(&newVideo.ID); err != nil {
		errors.DebugPrintf(err)
		if errors.IsForeignKeyError(err) {
			return &model.Video{}, errors.UserNotExist
		}
		return &model.Video{}, errors.InternalServerError
	}

	return newVideo, nil
}

// Videos is the resolver for the Videos field.
func (r *queryResolver) Videos(ctx context.Context, limit *int, offset *int) ([]*model.Video, error) {
	var video *model.Video
	var videos []*model.Video

	rows, err := dal.LogAndQuery(r.db, "SELECT id, name, url, created_at, user_id FROM videos ORDER BY created_at desc limit $1 offset $2", limit, offset)
	defer rows.Close()

	if err != nil {
		errors.DebugPrintf(err)
		return nil, errors.InternalServerError
	}
	for rows.Next() {
		if err := rows.Scan(&video.ID, &video.Name, &video.URL, &video.CreatedAt, video.User.ID); err != nil {
			errors.DebugPrintf(err)
			return nil, errors.InternalServerError
		}
		videos = append(videos, video)
	}

	return videos, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
