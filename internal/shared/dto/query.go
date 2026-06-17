package dto

import (
	"strings"

	"github.com/XaiPhyr/rdev-go-api-template/internal/shared/helpers"
)

type Query struct {
	Limit  int    `json:"limit" query:"limit"`
	Offset int    `json:"offset" query:"offset"`
	Order  string `json:"order" query:"order"`
	Search string `json:"search" query:"search"`
}

type BaseFilters struct {
	Page     int    `form:"page,default=0"`
	PageSize int    `form:"page_size,default=10"`
	Search   string `form:"search" binding:"omitempty,alphanumspace"`
	Sort     string `form:"sort" binding:"omitempty,alphaspace"`
}

func (b BaseFilters) SanitizeQuery(allowedColumns []string) Query {
	pageSize := b.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	if b.Search != "" {
		trimmedSpaces := strings.TrimSpace(b.Search)
		b.Search = helpers.CleanSpecialChars(trimmedSpaces)
	}

	return Query{
		Limit:  pageSize,
		Offset: (b.Page - 1) * pageSize,
		Order:  b.validateSort(allowedColumns),
		Search: b.Search,
	}
}

func (b BaseFilters) validateSort(allowedColumns []string) string {
	finalSort := "id ASC"
	if b.Sort != "" {
		for _, col := range allowedColumns {
			if b.Sort == col || b.Sort == col+" ASC" || b.Sort == col+" DESC" {
				finalSort = b.Sort
				break
			}
		}
	}

	return finalSort
}
