package shared

import (
	"context"

	"github.com/startdusk/twitter/data"
)

type UserIDKey struct{}

func GetUserIDFromContext(c context.Context) (string, error) {
	val := c.Value(UserIDKey{})
	userID, ok := val.(string)
	if !ok {
		return "", data.ErrNoUserIDInContext
	}
	return userID, nil
}
