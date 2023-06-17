package fabric

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type Pagination struct {
	Page      int    `query:"page"`
	Size      int    `query:"size"`
	OrderBy   string `query:"order"`
	OrderDesc bool   `query:"desc"`
}

func (obj Pagination) GetOffset(page, size int) int {
	return (obj.Page - 1) * obj.Size
}

func (obj Pagination) GormPaginate(tx *gorm.DB) *gorm.DB {
	page := obj.GetOffset(DEFAULT_PAGE, DEFAULT_PAGE_SIZE)
	size := DEFAULT_PAGE_SIZE

	if obj.Size > 0 {
		size = obj.Size
	}

	if obj.Page > 0 {
		page = obj.GetOffset(obj.Page, size)
	}

	return tx.Offset(page).Limit(size)
}

func (obj Pagination) GormOrder(tx *gorm.DB, defaultOrderBy string) *gorm.DB {
	order := fmt.Sprintf("%v ASC", defaultOrderBy)
	if obj.OrderBy != "" {
		order = fmt.Sprintf("%v ASC", defaultOrderBy)
	}

	if obj.OrderDesc {
		order = strings.ReplaceAll(order, "ASC", "DESC")
	}

	return tx.Order(order)
}
