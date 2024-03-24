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

type GetFriendListRequest struct {
	UserId     string `json:"userId"`
	Limit      int    `form:"limit" query:"limit" json:"limit"`
	Offset     int    `form:"offset" query:"offset" json:"offset"`
	SortBy     string `form:"sortBy" query:"sortBy" json:"sortBy"`
	OrderBy    string `form:"orderBy" query:"orderBy" json:"orderBy"`
	OnlyFriend bool   `form:"onlyFriend" query:"onlyFriend" json:"onlyFriend"`
	Search     string `form:"search" query:"search" json:"search"`
}

type FriendResponse struct {
	UserId      string    `json:"userId"`
	FriendId    string    `json:"friendId,omitempty"`
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

func (p GetFriendListRequest) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Limit, validation.Min(0), validation.Max(100)),
		validation.Field(&p.Offset, validation.Min(0)),
		validation.Field(&p.OrderBy, validation.In("asc", "desc")),
		validation.Field(&p.SortBy, validation.In("friendCount", "createdAt")),
		// validation.Field(&p.OnlyFriend, validation.Min(0), validation.Max(1)),
	)
}
