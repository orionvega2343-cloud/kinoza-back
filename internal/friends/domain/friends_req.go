package domain

import "time"

const (
	StatusPending  = "pending"
	StatusAccepted = "accepted"
	StatusDeclined = "declined"
)

type FriendRequest struct {
	ID         int       `db:"id"`
	FromUserID string    `db:"from_user_id"`
	ToUserID   string    `db:"to_user_id"`
	Status     string    `db:"status"`
	CreatedAt  time.Time `db:"created_at"`
}

func (r *FriendRequest) Validate() error {
	if r.FromUserID == r.ToUserID {
		return ErrSelfRequest
	}
	return nil
}

func (r *FriendRequest) Accept(userID string) error {
	if err := r.canRespond(userID); err != nil {
		return err
	}
	r.Status = StatusAccepted
	return nil
}

func (r *FriendRequest) Decline(userID string) error {
	if err := r.canRespond(userID); err != nil {
		return err
	}
	r.Status = StatusDeclined
	return nil
}

func (r *FriendRequest) canRespond(userID string) error {
	if r.ToUserID != userID {
		return ErrNotRecipient
	}
	if r.Status != StatusPending {
		return ErrInvalidTransition
	}
	return nil
}
