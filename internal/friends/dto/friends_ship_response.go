package dto

import "time"

type FriendResponse struct {
	ID        int       `json:"id"`
	FriendID  string    `json:"friend_id"`
	CreatedAt time.Time `json:"created_at"`
}
