package response

import (
	"time"

	"github.com/Delkira544/rakiduam/internal/shared/errors"
)

type APIResponse struct {
	Success bool      `json:"success"`
	Data    any       `json:"data,omitempty"`
	Error   *APIError `json:"error,omitempty"`
	Meta    *Meta     `json:"meta,omitempty"`
}

type APIError struct {
	Code    string              `json:"code"`
	Message string              `json:"message"`
	Details []errors.FieldError `json:"details,omitempty"`
}

type Meta struct {
	RequestID  string `json:"request_id"`
	Timestamp  string `json:"timestamp"`
	Page       int    `json:"page,omitempty"`
	PerPage    int    `json:"per_page,omitempty"`
	Total      int64  `json:"total,omitempty"`
	TotalPages int    `json:"total_pages,omitempty"`
	HasNext    bool   `json:"has_next,omitempty"`
	HasPrev    bool   `json:"has_prev,omitempty"`
}

func NewSuccessResponse(data any) APIResponse {
	return APIResponse{Success: true, Data: data}
}

func NewSuccessResponseWithMeta(data any, meta *Meta) APIResponse {
	return APIResponse{Success: true, Data: data, Meta: meta}
}

func NewErrorResponse(appErr *errors.AppError) APIResponse {
	return APIResponse{
		Success: false,
		Error: &APIError{
			Code:    appErr.Code,
			Message: appErr.Message,
			Details: appErr.Details,
		},
	}
}

func NewMeta() *Meta {
	return &Meta{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

func NewPaginatedMeta(page, perPage int, total int64) *Meta {
	meta := NewMeta()
	meta.Page = page
	meta.PerPage = perPage
	meta.Total = total
	if perPage > 0 {
		meta.TotalPages = int((total + int64(perPage) - 1) / int64(perPage))
	}
	meta.HasNext = page < meta.TotalPages
	meta.HasPrev = page > 1
	return meta
}
