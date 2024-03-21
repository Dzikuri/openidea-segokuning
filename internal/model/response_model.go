package model

type PaginateResponse[T any] struct {
	Data    []T              `json:"data"`
	Meta    MetaDataResponse `json:"meta,omitempty"`
	Errors  string           `json:"errors,omitempty"`
	Message string           `json:"message"`
}

type Response[T any] struct {
	Code    int    `json:"-"`
	Data    T      `json:"data"`
	Message string `json:"message"`
}

type ResponseError struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
	Error   error  `json:"error,omitempty"`
}

type MetaDataResponse struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}
