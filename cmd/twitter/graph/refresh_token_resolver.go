package graph

import "context"

func (m *mutationResolver) RefreshToken(ctx context.Context, input RefreshTokenInput) (*AuthResponse, error) {
	resp, err := m.AuthService.RefreshToken(ctx, input.Token)
	if err != nil {
		return nil, err
	}
	return &AuthResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		User:         &User{ID: resp.User.ID},
	}, err
}
