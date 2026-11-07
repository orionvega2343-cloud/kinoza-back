package domain

import "errors"

var (
	ErrRequestNotFound      = errors.New("friend request not found")
	ErrSelfRequest          = errors.New("cannot send friend request to yourself")
	ErrRequestAlreadyExists = errors.New("pending friend request already exists")
	ErrAlreadyFriends       = errors.New("users are already friends")
	ErrNotRecipient         = errors.New("only the recipient can respond to the request")
	ErrInvalidTransition    = errors.New("friend request is not pending")
	ErrFriendshipNotFound   = errors.New("friendship not found")
	ErrUserNotFound         = errors.New("user not found")
)
