package sql

import (
	"github.com/70X/properties-seeker/storage"
	"github.com/70X/properties-seeker/storage/src/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SQLProperties struct {
	conn *storage.DB
}

func NewSQLProperties(conn *storage.DB) *SQLProperties {
	return &SQLProperties{conn: conn}
}

func (s SQLProperties) queryByContractType(contractType models.ContractType) *gorm.DB {
	query := s.conn.Preload("Execution", "contract_type = ?", contractType).
		Joins("JOIN executions E ON properties.execution_id = E.id").
		Where("E.contract_type = ?", contractType)
	return query
}

func (s SQLProperties) whereArchived(archived bool) string {
	if !archived {
		return "archived_at IS NULL"
	}
	return "archived_at IS NOT NULL"
}

func (s SQLProperties) GetProperty(id uint) (*models.Property, error) {
	item := models.Property{Model: gorm.Model{ID: id}}
	result := s.conn.Preload("Execution").First(&item)
	if result.Error != nil {
		return nil, result.Error
	}
	return &item, nil
}

func (s SQLProperties) FetchRentProperties(archived bool) ([]models.Property, error) {
	properties := make([]models.Property, 0)
	query := s.queryByContractType(models.RENT)
	query = query.Where(s.whereArchived(archived))

	err := query.Find(&properties).Error

	return properties, err
}

func (s SQLProperties) FetchSellProperties(archived bool) ([]models.Property, error) {
	properties := make([]models.Property, 0)
	query := s.queryByContractType(models.SELL).Where(models.Property{IsAuction: false})
	query = query.Where(s.whereArchived(archived))

	err := query.Find(&properties).Error

	return properties, err
}

func (s SQLProperties) FetchAuctionList(archived bool) ([]models.Property, error) {
	properties := make([]models.Property, 0)
	query := s.conn.Where(models.Property{IsAuction: true})
	query = query.Where(s.whereArchived(archived))

	err := query.Find(&properties).Error
	return properties, err
}

func (s SQLProperties) AddProperties(properties []models.Property) ([]models.Property, error) {
	result := s.conn.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "platform_id"}},
		DoNothing: true,
	}).Create(&properties)

	if result.Error != nil {
		return nil, result.Error
	}
	items := make([]models.Property, 0)
	for _, property := range properties {
		if property.ID != 0 {
			items = append(items, property)
		}
	}
	return items, nil
}

func (s SQLProperties) UpdateProperty(property models.Property) (*models.Property, error) {
	res := s.conn.Save(&property)
	if res.Error != nil {
		return nil, res.Error
	}
	return &property, nil
}

func (s SQLProperties) DeleteProperty(id uint) error {
	err := s.conn.Delete(&models.Property{}, id).Error
	return err
}
