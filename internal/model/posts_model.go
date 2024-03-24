package model

import (
	"time"

	validation "github.com/itgelo/ozzo-validation/v4"
	"github.com/lib/pq"
)

type CreatePostRequest struct {
	UserId  string   `json:"userId"`
	Content string   `json:"postInHtml"`
	Tags    []string `json:"tags,omitempty"`
}

type CreatePostCommentRequest struct {
	PostId  string `json:"postId"`
	UserId  string `json:"userId"`
	Comment string `json:"comment"`
}

type PostResponse struct {
	Id        string         `json:"id"`
	UserId    string         `json:"userId"`
	Content   string         `json:"postInHtml"`
	Tags      pq.StringArray `json:"tags"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
}

type PostCommentResponse struct {
	Id        string    `json:"id"`
	PostId    string    `json:"postId"`
	UserId    string    `json:"userId"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type PostCommentUserResponse struct {
	PostCommentResponse
	Creator FriendResponse `json:"creator"`
}

type PostListResponse struct {
	PostId   string                    `json:"postId"`
	Post     PostResponse              `json:"post"`
	Comments []PostCommentUserResponse `json:"comments"`
	Creator  FriendResponse            `json:"creator"`
}

type PostListRequest struct {
	UserId string `json:"userId"`
	Limit  int    `form:"limit" query:"limit, omitempty" json:"limit"`
	Offset int    `form:"offset" query:"offset, omitempty" json:"offset"`
	// SortBy     string `form:"sortBy" query:"sortBy" json:"sortBy"`
	// OrderBy    string `form:"orderBy" query:"orderBy" json:"orderBy"`
	Search    string   `form:"search" query:"search" json:"search"`
	SearchTag []string `form:"searchTag" query:"searchTag" json:"searchTag"`
}

func (r CreatePostRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Content, validation.Required, validation.Length(2, 500)),
		validation.Field(&r.Tags, validation.Required, validation.Each(validation.Required), validation.Length(1, 20)),
	)
}

func (r CreatePostCommentRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Comment, validation.Required, validation.Length(2, 500)),
		validation.Field(&r.PostId, validation.Required),
	)
}

func (p PostListRequest) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Limit, validation.Min(0), validation.Max(100), validation.When(p.Limit != 0, validation.Required)),
		validation.Field(&p.Offset, validation.Min(0), validation.When(p.Offset != 0, validation.Required)),
	)
}
