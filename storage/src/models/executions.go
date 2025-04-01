package models

import (
	"gorm.io/gorm"
)

type ExecutionStatus string

const (
	SUCCESS  ExecutionStatus = "SUCCESS"
	FAILURE  ExecutionStatus = "FAILURE"
	PROGRESS ExecutionStatus = "IN_PROGRESS"
	WARNING  ExecutionStatus = "WARNING"
)

type Execution struct {
	gorm.Model
	ScraperFilter `gorm:"embedded"`
	Status        ExecutionStatus
	Total         int
	Results       int
	Properties    []Property `gorm:"foreignKey:ExecutionID"`
}
