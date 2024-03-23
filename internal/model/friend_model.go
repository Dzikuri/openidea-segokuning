package model

import (
	"time"

	validation "github.com/itgelo/ozzo-validation/v4"
	"github.com/itgelo/ozzo-validation/v4/is"
)

type FriendRequest struct {
	UserId   string `json:"userId"`
	FriendId string `json:"friendId"`
}

type FriendResponse struct {
	UserId      string    `json:"userId"`
	FriendId    string    `json:"friendId"`
	Name        string    `json:"name"`
	ImageUrl    string    `json:"imageUrl"`
	FriendCount int       `json:"friendCount"`
	CreatedAt   time.Time `json:"createdAt"`
}

func (p FriendRequest) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.UserId, validation.Required, is.UUIDv4),
	)
}
