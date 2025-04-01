package models

import "gorm.io/gorm"

type ContractType string

const (
	RENT ContractType = "RENT"
	SELL ContractType = "SELL"
)

type PlatformSource string

const (
	IMMOBILIARE PlatformSource = "IMMOBILIARE"
	IDEALISTA   PlatformSource = "IDEALISTA"
	CASA        PlatformSource = "CASA"
)

type ScraperFilter struct {
	gorm.Model
	ContractType   ContractType   `json:"contract_type"`
	Area           string         `json:"area"`
	PlatformSource PlatformSource `json:"platform_source"`
	FromSize       int            `json:"from_size"`
	MaxPrice       int            `json:"max_price"`
	FromPage       int            `json:"from_page"`
	ToPage         int            `json:"to_page"`
	Enabled        bool           `json:"enabled"`
}
