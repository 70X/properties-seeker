package sql

import (
	"github.com/70X/properties-seeker/storage"
	"github.com/70X/properties-seeker/storage/src/models"
	"gorm.io/gorm"
)

type SQLExecutions struct {
	conn *storage.DB
}

func NewSQLExecutions(conn *storage.DB) *SQLExecutions {
	return &SQLExecutions{conn}
}

func (r *SQLExecutions) StartExecution(filter models.ScraperFilter) (models.Execution, error) {
	item := models.Execution{ScraperFilter: filter, Status: models.PROGRESS}
	res := r.conn.Create(&item)
	return item, res.Error
}

func (r *SQLExecutions) SetExecutionStatus(id uint, status models.ExecutionStatus, total int, results int) (models.Execution, error) {
	item := models.Execution{Model: gorm.Model{ID: id}}
	r.conn.First(&item)
	item.Status = status
	item.Total = total
	item.Results = results
	res := r.conn.Save(&item)
	return item, res.Error
}

func (r *SQLExecutions) GetExecutionByFilter(filter models.ScraperFilter) (models.Execution, error) {
	item := models.Execution{ScraperFilter: filter}
	res := r.conn.Order("updated_at desc").Where(
		"area = ? AND platform_source = ? AND contract_type = ? AND from_page = ? AND to_page = ? AND from_size = ?",
		filter.Area, filter.PlatformSource, filter.ContractType, filter.FromPage, filter.ToPage, filter.FromSize,
	).First(&item)
	if item.ID == 0 {
		return models.Execution{}, res.Error
	}
	return item, res.Error
}
