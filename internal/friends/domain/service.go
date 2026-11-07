package domain

import "context"

type FriendRequestService interface {
	Send(ctx context.Context, r *FriendRequest) error
	Accept(ctx context.Context, requestID int, userID string) error
	Decline(ctx context.Context, requestID int, userID string) error
	ListIncoming(ctx context.Context, userID string) ([]FriendRequest, error)
	ListOutgoing(ctx context.Context, userID string) ([]FriendRequest, error)
}

type FriendshipService interface {
	List(ctx context.Context, userID string) ([]Friendship, error)
	Delete(ctx context.Context, userID, friendID string) error
}
