//go:generate go run github.com/vektah/dataloaden UserLoader string *github.com/startdusk/twitter/cmd/twitter/graph.User
package graph

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/startdusk/twitter/data"
)

type loadersKey struct{}

type Loaders struct {
	UserByID UserLoader
}

type Repos struct {
	UserRepo data.UserRepo
}

func DataloaderMiddleware(repos *Repos) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		c.Request = c.Request.WithContext(context.WithValue(ctx, loadersKey{}, &Loaders{
			UserByID: UserLoader{
				fetch: func(ids []string) ([]*User, []error) {
					users := make([]*User, len(ids))
					var errors []error
					repoUsers, err := repos.UserRepo.GetByIDs(ctx, ids)
					if err != nil {
						errors = append(errors, err)
						return nil, errors
					}
					userByID := make(map[string]*User)
					for _, u := range repoUsers {
						userByID[u.ID] = mapToUser(u)
					}
					for i, id := range ids {
						user, ok := userByID[id]
						if !ok {
							errors = append(errors, fmt.Errorf("user with id[%s] is missing", id))
						}
						users[i] = user
					}

					return users, errors
				},
				wait:     1 * time.Millisecond,
				maxBatch: 100,
				batch: &userLoaderBatch{
					closing: false,
					done:    make(chan struct{}),
				},
			},
		}))
		c.Next()
	}
}

func DataloaderFor(ctx context.Context) *Loaders {
	return ctx.Value(loadersKey{}).(*Loaders)
}
