package service

import (
	"context"

	"kinoza-back/internal/friends/domain"
)

type FriendshipService struct {
	repo domain.FriendshipRepository
}

func NewFriendshipService(repo domain.FriendshipRepository) *FriendshipService {
	return &FriendshipService{
		repo: repo,
	}
}

var _ domain.FriendshipService = (*FriendshipService)(nil)

func (s *FriendshipService) List(ctx context.Context, userID string) ([]domain.Friendship, error) {
	return s.repo.List(ctx, userID)
}

func (s *FriendshipService) Delete(ctx context.Context, userID, friendID string) error {
	return s.repo.Delete(ctx, domain.NewFriendship(userID, friendID))
}
