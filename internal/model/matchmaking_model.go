package model

type CreateMatchMakingRequest struct {
	UserID uint64 `json:"user_id" validate:"required"`
}

type CreateMatchMakingResponse struct {
	RoomID string   `json:"room_id"`
	Users  []string `json:"users"`
}
