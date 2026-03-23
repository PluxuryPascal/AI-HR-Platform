package domain

type Pagination struct {
	Page    int `json:"page" validate:"required,min=1"`
	PerPage int `json:"per_page" validate:"required,min=1,max=1000"`
}
