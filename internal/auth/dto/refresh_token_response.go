package dto

import "time"

type RefreshTokenResponse struct {
	Id        string    `json:"id"`
	UserId    string    `json:"user_id"`
	TokenHash string    `json:"token_hash"`
	ExpiresAt time.Time `json:"expires_at"`
	RevokedAt time.Time `json:"revoked_at"`
	CreatedAt time.Time `json:"created_at"`
}
