package service

import (
	"context"

	"kinoza-back/internal/friends/domain"
	"kinoza-back/pkg/transaction"
)

type FriendsService struct {
	requests    domain.FriendRequestRepository
	friendships domain.FriendshipRepository
	tx          *transaction.Transactor
}

func NewFriendsService(requests domain.FriendRequestRepository, friendships domain.FriendshipRepository, tx *transaction.Transactor) *FriendsService {
	return &FriendsService{
		requests:    requests,
		friendships: friendships,
		tx:          tx,
	}
}

var _ domain.FriendRequestService = (*FriendsService)(nil)

func (s *FriendsService) Send(ctx context.Context, req *domain.FriendRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}
	exists, err := s.friendships.Exists(ctx, domain.NewFriendship(req.FromUserID, req.ToUserID))
	if err != nil {
		return err
	}
	if exists {
		return domain.ErrAlreadyFriends
	}
	return s.requests.Create(ctx, req)
}

func (s *FriendsService) ListIncoming(ctx context.Context, userID string) ([]domain.FriendRequest, error) {
	return s.requests.ListIncoming(ctx, userID)
}

func (s *FriendsService) ListOutgoing(ctx context.Context, userID string) ([]domain.FriendRequest, error) {
	return s.requests.ListOutgoing(ctx, userID)
}

func (s *FriendsService) Decline(ctx context.Context, requestID int, userID string) error {
	return s.tx.Transaction(ctx, func(ctx context.Context) error {
		req, err := s.requests.GetByIDForUpdate(ctx, requestID)
		if err != nil {
			return err
		}
		if err := req.Decline(userID); err != nil {
			return err
		}
		return s.requests.UpdateStatus(ctx, req)
	})
}

func (s *FriendsService) Accept(ctx context.Context, requestID int, userID string) error {
	return s.tx.Transaction(ctx, func(ctx context.Context) error {
		req, err := s.requests.GetByIDForUpdate(ctx, requestID)
		if err != nil {
			return err
		}
		if err := req.Accept(userID); err != nil {
			return err
		}
		if err := s.requests.UpdateStatus(ctx, req); err != nil {
			return err
		}
		return s.friendships.Create(ctx, domain.NewFriendship(req.FromUserID, req.ToUserID))
	})
}
