package util

import (
	"net/http"
	"strconv"
)

type Pagination struct {
	Page   int `json:"page"`
	Limit  int `json:"limit"`
	Offset int `json:"-"` // db logic offset
}

// NewPagination extracts and validates query parameters
func NewPagination(r *http.Request) Pagination {
	query := r.URL.Query()

	page, err := strconv.Atoi(query.Get("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err := strconv.Atoi(query.Get("limit"))
	if err != nil || limit <= 0 {
		limit = 20
	} else if limit > 100 {
		limit = 100
	}

	return Pagination{
		Page:   page,
		Limit:  limit,
		Offset: (page - 1) * limit,
	}
}
