package sql

import (
	"github.com/70X/properties-seeker/storage"
	"github.com/70X/properties-seeker/storage/src/models"
	"gorm.io/gorm"
)

type SQLScraperFilter struct {
	conn *storage.DB
}

func NewSQLScraperFilter(conn *storage.DB) *SQLScraperFilter {
	return &SQLScraperFilter{conn: conn}
}

func (r *SQLScraperFilter) GetScraperFilters() ([]models.ScraperFilter, error) {
	items := make([]models.ScraperFilter, 0)
	res := r.conn.Where("enabled = ?", true).Find(&items)
	if res.Error != nil {
		return nil, res.Error
	}
	return items, nil
}

func (r *SQLScraperFilter) AddScraperFilter(filter models.ScraperFilter) (uint, error) {
	item := filter
	item.Enabled = true
	res := r.conn.Create(&item)
	return item.ID, res.Error
}

func (r *SQLScraperFilter) EnableScraperFilter(id uint, enable bool) (uint, error) {
	item := models.ScraperFilter{Model: gorm.Model{ID: id}}
	r.conn.First(&item)
	item.Enabled = enable
	res := r.conn.Save(&item)
	return item.ID, res.Error
}
