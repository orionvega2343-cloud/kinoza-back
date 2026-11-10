package handler

import (
	"kinoza-back/internal/friends/domain"
	"kinoza-back/internal/friends/dto"
)

func toFriendRequestResponse(r *domain.FriendRequest) dto.FriendRequestResponse {
	return dto.FriendRequestResponse{
		ID:         r.ID,
		FromUserID: r.FromUserID,
		ToUserID:   r.ToUserID,
		Status:     r.Status,
		CreatedAt:  r.CreatedAt,
	}
}

func toFriendRequestListResponse(requests []domain.FriendRequest) []dto.FriendRequestResponse {
	res := make([]dto.FriendRequestResponse, 0, len(requests))

	for _, r := range requests {
		res = append(res, toFriendRequestResponse(&r))
	}
	return res
}

// toFriendListResponse — userID нужен, чтобы из пары (a, b) отдать id именно друга.
func toFriendListResponse(friendships []domain.Friendship, userID string) []dto.FriendResponse {
	res := make([]dto.FriendResponse, 0, len(friendships))

	for _, f := range friendships {
		res = append(res, dto.FriendResponse{
			ID:        f.ID,
			FriendID:  f.FriendOf(userID),
			CreatedAt: f.CreatedAt,
		})
	}
	return res
}
