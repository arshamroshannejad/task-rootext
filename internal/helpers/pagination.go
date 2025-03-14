package helpers

import (
	"math"
	"strings"
)

type PaginateFilter struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafeList []string
}

func (p *PaginateFilter) Validate(v *Validator) {
	v.Check(p.Page > 0, "page", "must be greater than zero")
	v.Check(p.Page <= 10_000_000, "page", "must be a maximum of 10 million")
	v.Check(p.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(p.PageSize <= 100, "page_size", "must be a maximum of 100")
	v.Check(v.In(p.Sort, p.SortSafeList...), "sort", "invalid sort value")
}

func (p *PaginateFilter) Limit() int {
	return p.PageSize
}

func (p *PaginateFilter) OffSet() int {
	return (p.Page - 1) * p.PageSize
}

func (p *PaginateFilter) SortValue() string {
	for _, v := range p.SortSafeList {
		if v == p.Sort {
			return strings.TrimPrefix(p.Sort, "-")
		}
	}
	panic("unsafe sort parameter: " + p.Sort)
}

func (p *PaginateFilter) SortDirection() string {
	if strings.HasPrefix(p.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}

type Metadata struct {
	CurrentPage  int `json:"current_page,omitempty" example:"3"`
	PageSize     int `json:"page_size,omitempty" example:"2"`
	FirstPage    int `json:"first_page,omitempty" example:"1"`
	LastPage     int `json:"last_page,omitempty" example:"12"`
	TotalRecords int `json:"total_records,omitempty" example:"1200"`
}

func CalculateMetadata(totalRecords, page, pageSize int) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}
	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}
}
