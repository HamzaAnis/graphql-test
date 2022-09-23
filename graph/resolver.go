package graph

import (
	"context"
	"database/sql"
	"github.com/99designs/gqlgen/graphql"
	"github.com/hamzaanis/graphql-test/graph/errors"
	"github.com/hamzaanis/graphql-test/graph/generated"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type contextKey string

var (
	UserIDCtxKey = contextKey("userID")
)

type Resolver struct {
	db *sql.DB
}

func NewRootResolvers(db *sql.DB) generated.Config {
	c := generated.Config{
		Resolvers: &Resolver{
			db: db,
		},
	}

	// Complexity
	countComplexity := func(childComplexity int, limit *int, offset *int) int {
		return *limit * childComplexity
	}
	c.Complexity.Query.Videos = countComplexity
	c.Complexity.Video.Related = countComplexity

	// Schema Directive
	c.Directives.IsAuthenticated = func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
		ctxUserID := ctx.Value(UserIDCtxKey)
		if ctxUserID != nil {
			return next(ctx)
		} else {
			return nil, errors.UnauthorisedError
		}
	}
	return c
}
