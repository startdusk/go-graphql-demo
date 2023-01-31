package jwt

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/startdusk/twitter/config"
	"github.com/startdusk/twitter/data"
	"github.com/stretchr/testify/assert"
)

var (
	authTokenService data.AuthTokenService
)

func testNowFunc() time.Time {
	// 固定时间为 2023-01-31 12:21:06
	return time.Date(2023, time.January, 31, 12, 21, 6, 0, time.Local)
}

func TestMain(m *testing.M) {
	nowFunc = testNowFunc
	conf, err := config.New("../.env.test")
	if err != nil {
		panic(err)
	}

	authTokenService = NewTokenService(&conf.JWT)
	m.Run()
}

func TestTokenService_CreateAndParseAccessToken(t *testing.T) {
	t.Run("can create access token", func(t *testing.T) {
		ctx := context.Background()
		accessToken, err := authTokenService.CreateAccessToken(ctx, data.User{
			ID: "user_id",
		})
		assert.NoError(t, err)

		token, err := authTokenService.ParseToken(ctx, accessToken)
		assert.NoError(t, err)
		assert.Empty(t, token.ID)
		assert.Equal(t, token.Sub, token.Sub)
	})

	t.Run("return err if invalid access token", func(t *testing.T) {
		token, err := authTokenService.ParseToken(context.Background(), "invalid token")
		assert.Error(t, err)
		assert.Equal(t, token, data.NilAuthToken)
	})

	t.Run("expired time", func(t *testing.T) {
		ctx := context.Background()
		accessToken, err := authTokenService.CreateAccessToken(ctx, data.User{
			ID: "user_id",
		})
		assert.NoError(t, err)
		nowFunc = func() time.Time {
			return time.Now().Add(data.AccessTokenLifetime)
		}
		token, err := authTokenService.ParseToken(ctx, accessToken)
		assert.Error(t, err)
		assert.Equal(t, token, data.NilAuthToken)
		resetNowFunc(t)
	})
}

func TestTokenService_CreateAndParseRefreshToken(t *testing.T) {
	t.Run("can create refresh token", func(t *testing.T) {
		ctx := context.Background()
		tokenID := "token_id"
		refreshToken, err := authTokenService.CreateRefreshToken(ctx, data.User{
			ID: "user_id",
		}, tokenID)
		assert.NoError(t, err)

		token, err := authTokenService.ParseToken(ctx, refreshToken)
		assert.NoError(t, err)
		assert.Equal(t, token.ID, tokenID)
		assert.Equal(t, token.Sub, token.Sub)
	})
}

func TestTokenService_ParseFromRequest(t *testing.T) {
	t.Run("can parse access token from request", func(t *testing.T) {
		ctx := context.Background()
		accessToken, err := authTokenService.CreateAccessToken(ctx, data.User{
			ID: "user_id",
		})
		assert.NoError(t, err)

		req, err := http.NewRequest(http.MethodGet, `https://twitter.clone.com`, nil)
		assert.NoError(t, err)
		req.Header.Set(`Authorization`, fmt.Sprintf(`Bearer %s`, accessToken))

		token, err := authTokenService.ParseTokenFromRequest(ctx, req)
		assert.NoError(t, err)

		assert.Empty(t, token.ID)
		assert.Equal(t, token.Sub, token.Sub)
	})

	t.Run("can parse refresh token from request", func(t *testing.T) {
		ctx := context.Background()
		tokenID := "token_id"
		refreshToken, err := authTokenService.CreateRefreshToken(ctx, data.User{
			ID: "user_id",
		}, tokenID)
		assert.NoError(t, err)

		req, err := http.NewRequest(http.MethodGet, `https://twitter.clone.com`, nil)
		assert.NoError(t, err)
		req.Header.Set(`Authorization`, fmt.Sprintf(`Bearer %s`, refreshToken))

		token, err := authTokenService.ParseTokenFromRequest(ctx, req)
		assert.NoError(t, err)

		assert.Equal(t, token.ID, tokenID)
		assert.Equal(t, token.Sub, token.Sub)
	})
}

func resetNowFunc(t *testing.T) {
	t.Helper()

	nowFunc = testNowFunc
}
