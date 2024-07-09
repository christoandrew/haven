package pagination

import (
	"math"

	"gorm.io/gorm"
)

type Response struct {
	TotalCount int         `json:"total_count"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	PrevPage   int         `json:"prev_page"`
	NextPage   int         `json:"next_page"`
	LastPage   int         `json:"last_page"`
	Results    interface{} `json:"results"`
}

type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalCount int `json:"total_count"`
	serializer interface{}
}

func (pagination *Pagination) PrevPage() int {
	if pagination.Page > 1 {
		return pagination.Page - 1
	}

	return 1
}

func (pagination *Pagination) NextPage() int {
	if pagination.TotalCount > pagination.Page*pagination.Limit {
		return pagination.Page + 1
	}

	return pagination.TotalCount / pagination.Limit
}

func (pagination *Pagination) LastPage() int {
	return int(math.Ceil(float64(pagination.TotalCount) / float64(pagination.Limit)))
}

func (pagination *Pagination) Paginate(db *gorm.DB, model interface{}) *gorm.DB {
	var totalCount int64
	db.Model(model).Count(&totalCount)
	pagination.TotalCount = int(totalCount)
	return db.Limit(pagination.Limit).Offset((pagination.Page - 1) * pagination.Limit).Order("id asc")
}
