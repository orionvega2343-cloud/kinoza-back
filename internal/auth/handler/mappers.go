package handler

import (
	"kinoza-back/internal/auth/domain"
	"kinoza-back/internal/auth/dto"
)

func toUserDomainUser(req dto.UserRequest) domain.User {
	return domain.User{
		Email:        req.Email,
		PasswordHash: req.Password,
	}
}

func toUserResponse(user domain.User) dto.UserResponse {
	return dto.UserResponse{
		Id:           user.Id,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Role:         user.Role,
		CreatedAt:    user.CreatedAt,
	}
}

func toRefreshTokenResponse(rt *domain.RefreshToken) dto.RefreshTokenResponse {
	return dto.RefreshTokenResponse{
		Id:        rt.Id,
		UserId:    rt.UserId,
		TokenHash: rt.TokenHash,
		ExpiresAt: rt.ExpiresAt,
		RevokedAt: rt.RevokedAt,
		CreatedAt: rt.CreatedAt,
	}
}
