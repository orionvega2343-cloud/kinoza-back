package dto

type SendFriendRequest struct {
	ToUserID string `json:"to_user_id" binding:"required,uuid"`
}
