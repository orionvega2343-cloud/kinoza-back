package domain

import "context"

type FriendRequestRepository interface {
	Create(ctx context.Context, r *FriendRequest) error
	// GetByIDForUpdate блокирует строку (SELECT ... FOR UPDATE) — вызывать только внутри транзакции.
	GetByIDForUpdate(ctx context.Context, id int) (*FriendRequest, error)
	UpdateStatus(ctx context.Context, r *FriendRequest) error
	ListIncoming(ctx context.Context, userID string) ([]FriendRequest, error)
	ListOutgoing(ctx context.Context, userID string) ([]FriendRequest, error)
}

type FriendshipRepository interface {
	Create(ctx context.Context, f *Friendship) error
	Exists(ctx context.Context, f *Friendship) (bool, error)
	Delete(ctx context.Context, f *Friendship) error
	List(ctx context.Context, userID string) ([]Friendship, error)
}
