package jwt

import (
	"context"
	"net/http"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/pkg/errors"
	"github.com/startdusk/twitter/config"
	"github.com/startdusk/twitter/data"
)

var signatureType = jwa.HS256

var nowFunc = time.Now

type TokenService struct {
	Config *config.JWT
}

func NewTokenService(conf *config.JWT) *TokenService {
	return &TokenService{
		Config: conf,
	}
}

func (ts *TokenService) ParseTokenFromRequest(ctx context.Context, r *http.Request) (data.AuthToken, error) {
	token, err := jwt.ParseRequest(r,
		jwt.WithIssuer(ts.Config.Issuer),
		jwt.WithClock(jwt.ClockFunc(nowFunc)),
		jwt.WithKey(signatureType, []byte(ts.Config.Secret)),
	)
	if err != nil {
		return data.NilAuthToken, errors.Wrap(err, "parse token error from request")
	}
	return data.AuthToken{
		ID:  token.JwtID(),
		Sub: token.Subject(),
	}, nil
}

func (ts *TokenService) ParseToken(ctx context.Context, payload string) (data.AuthToken, error) {
	token, err := jwt.ParseString(payload,
		jwt.WithIssuer(ts.Config.Issuer),
		jwt.WithClock(jwt.ClockFunc(nowFunc)),
		jwt.WithKey(signatureType, []byte(ts.Config.Secret)),
	)
	if err != nil {
		return data.NilAuthToken, errors.Wrap(err, "parse token error from string")
	}

	return data.AuthToken{
		ID:  token.JwtID(),
		Sub: token.Subject(),
	}, nil
}

func (ts *TokenService) CreateRefreshToken(ctx context.Context, user data.User, tokenID string) (string, time.Time, error) {
	now := nowFunc()
	// Build a JWT!
	tok, err := jwt.NewBuilder().
		Subject(user.ID).
		Issuer(ts.Config.Issuer).
		IssuedAt(now).
		JwtID(tokenID).
		Expiration(now.Add(data.RefreshTokenLifetime)).
		Build()
	if err != nil {
		return "", time.Time{}, errors.Wrap(err, "failed to build refresh token")
	}

	// Sign a JWT!
	signed, err := jwt.Sign(tok, jwt.WithKey(signatureType, []byte(ts.Config.Secret)))
	if err != nil {
		return "", time.Time{}, errors.Wrap(err, "failed to sign refresh token")
	}
	return string(signed), tok.Expiration(), nil
}

func (ts *TokenService) CreateAccessToken(ctx context.Context, user data.User) (string, error) {
	now := nowFunc()
	// Build a JWT!
	tok, err := jwt.NewBuilder().
		Subject(user.ID).
		Issuer(ts.Config.Issuer).
		IssuedAt(now).
		Expiration(now.Add(data.AccessTokenLifetime)).
		Build()
	if err != nil {
		return "", errors.Wrap(err, "failed to build access token")
	}

	// Sign a JWT!
	signed, err := jwt.Sign(tok, jwt.WithKey(signatureType, []byte(ts.Config.Secret)))
	if err != nil {
		return "", errors.Wrap(err, "failed to sign access token")
	}

	return string(signed), nil
}
