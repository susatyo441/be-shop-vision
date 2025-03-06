package dto

type PaginationQuery struct {
	Page        int    `json:"page" transform:"int" validate:"gte=0"`
	Limit       int    `json:"limit" transform:"int" validate:"gte=0"`
	SortBy      string `json:"sortBy" transform:"string"`
	SortOrder   int    `json:"sortOrder" transform:"int"`
	Search      string `json:"search" transform:"string"`
	IsAvailable bool   `json:"isAvailable" transform:"bool"`
}

type ArrayOfIdDTO struct {
	IDs []string `json:"ids,omitempty" validate:"dive,mongodb"` // Optional list of category IDs
}
