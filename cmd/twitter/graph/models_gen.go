// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package graph

import (
	"time"
)

type AuthResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	User         *User  `json:"user"`
}

type CreatedTweetInput struct {
	Body string `json:"body"`
}

type LoginInput struct {
	UsernameOrEmail string `json:"usernameOrEmail"`
	Password        string `json:"password"`
}

type RefreshTokenInput struct {
	Token string `json:"token"`
}

type RegisterInput struct {
	Username        string `json:"username"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

type Tweet struct {
	ID        string    `json:"id"`
	Body      string    `json:"body"`
	User      *User     `json:"user"`
	UserID    string    `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
}

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}
