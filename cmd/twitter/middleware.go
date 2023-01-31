package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/startdusk/twitter/data"
	"github.com/startdusk/twitter/shared"
)

func authMiddleware(authTokenService data.AuthTokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		if token, err := authTokenService.ParseTokenFromRequest(ctx, c.Request); err == nil {
			c.Request = c.Request.WithContext(context.WithValue(ctx, shared.UserIDKey{}, token.Sub))
		}

		c.Next()
	}
}
