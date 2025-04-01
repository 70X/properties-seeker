package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Property struct {
	gorm.Model
	Title         string
	Price         int
	PriceText     string
	Agency        string
	ImageUrls     pq.StringArray `gorm:"type:text[]"`
	PageDetailUrl string
	Description   string
	PlatformID    string `gorm:"uniqueIndex"`
	Size          int
	Features      datatypes.JSON `gorm:"type:json"`
	OriginUrl     string
	Note          string
	ExecutionID   uint      `gorm:"index"`
	Execution     Execution `gorm:"foreignKey:ExecutionID"`
	ArchivedAt    *time.Time
	IsAuction     bool
}
