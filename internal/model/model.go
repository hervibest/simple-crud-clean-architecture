package model

type WebResponse struct {
	Success bool                   `json:"success"`
	Data    map[string]interface{} `json:"data,omitempty"`
	Paging  *PageMetadata          `json:"paging,omitempty"`
	Errors  string                 `json:"errors,omitempty"`
}

type ValidationErrorResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

type DataResponse[T any] struct {
	Success bool          `json:"success"`
	Data    T             `json:"data,omitempty"`
	Paging  *PageMetadata `json:"paging,omitempty"`
	Errors  string        `json:"errors,omitempty"`
}

// PageResponse remains the same
type PageResponse[T any] struct {
	Data         []T          `json:"data,omitempty"`
	PageMetadata PageMetadata `json:"paging,omitempty"`
}

type PageMetadata struct {
	Page            int    `json:"page"`
	Size            int    `json:"size"`
	Offset          int    `json:"offset,omitempty"`
	TotalItem       int64  `json:"total_item"`
	TotalPage       int64  `json:"total_page"`
	HasNext         bool   `json:"has_next"`
	HasPrevious     bool   `json:"has_previous"`
	NextPageURL     string `json:"next_page_url,omitempty"`
	PreviousPageURL string `json:"previous_page_url,omitempty"`
}
