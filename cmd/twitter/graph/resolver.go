//go:generate go run github.com/99designs/gqlgen generate

package graph

import (
	"context"
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	"github.com/startdusk/twitter/data"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type Resolver struct {
	AuthService  data.AuthService
	TweetService data.TweetService
	UserService  data.UserService
}

type queryResolver struct {
	*Resolver
}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct {
	*Resolver
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

type tweetResolver struct {
	*Resolver
}

func (r *Resolver) Tweet() TweetResolver {
	return &tweetResolver{r}
}

func writeBadRequestError(ctx context.Context, err error) error {
	return &gqlerror.Error{
		Message: err.Error(),
		Path:    graphql.GetPath(ctx),
		Extensions: map[string]any{
			"code": http.StatusBadRequest,
		},
	}
}

func writeUnauthenticatedError(ctx context.Context, err error) error {
	return &gqlerror.Error{
		Message: err.Error(),
		Path:    graphql.GetPath(ctx),
		Extensions: map[string]any{
			"code": http.StatusUnauthorized,
		},
	}
}
