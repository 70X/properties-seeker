package tests

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/70X/properties-seeker/storage"
	"github.com/70X/properties-seeker/storage/src/models"
	"gorm.io/datatypes"
)

func InitTestDb() *storage.DB {
	fmt.Println("Running SetupSuite (beforeAll)")
	config := &storage.Config{
		Host:     "localhost",
		Port:     "7444",
		User:     "user",
		Password: "password",
		DBName:   "properties_seeker",
	}
	conn := storage.NewDb(config)
	return conn
}

func BeforeEachTest(conn *storage.DB) {
	fmt.Println("Running SetupTest (beforeEach)")
	for _, model := range storage.MigrationModels {
		conn.Migrator().DropTable(model)
	}
	storage.RunMigrations(conn)
}

func mergeProperties(defaultProp, record models.Property) models.Property {
	defaultProp.PlatformID = GenerateRandomString(7)
	if record.PlatformID != "" {
		defaultProp.PlatformID = record.PlatformID
	}
	if record.ExecutionID != 0 {
		defaultProp.ExecutionID = record.ExecutionID
	}
	if record.Title != "" {
		defaultProp.Title = record.Title
	}
	if record.Price != 0 {
		defaultProp.Price = record.Price
	}
	if record.PriceText != "" {
		defaultProp.PriceText = record.PriceText
	}
	if record.Agency != "" {
		defaultProp.Agency = record.Agency
	}
	if len(record.ImageUrls) > 0 {
		defaultProp.ImageUrls = record.ImageUrls
	}
	if record.PageDetailUrl != "" {
		defaultProp.PageDetailUrl = record.PageDetailUrl
	}
	if record.Description != "" {
		defaultProp.Description = record.Description
	}
	if record.Size != 0 {
		defaultProp.Size = record.Size
	}
	if len(record.Features) > 0 {
		defaultProp.Features = record.Features
	}
	if record.OriginUrl != "" {
		defaultProp.OriginUrl = record.OriginUrl
	}
	if record.Note != "" {
		defaultProp.Note = record.Note
	}
	if record.ArchivedAt != nil {
		defaultProp.ArchivedAt = record.ArchivedAt
	}
	defaultProp.IsAuction = record.IsAuction
	return defaultProp
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateRandomString(n int) string {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	result := make([]byte, n)
	for i := range result {
		result[i] = charset[r.Intn(len(charset))]
	}

	return string(result)
}

func NewNProperties(record models.Property, n int) []models.Property {
	jsonStringValue := `{"Bedrooms": "2", "Bathrooms": "1", "Balcony": "Yes", "Parking": "Yes"}`
	defaultProperty := models.Property{
		ExecutionID:   1,
		Title:         "Cozy 2-Bedroom Apartment in Downtown",
		Price:         450000,
		Agency:        "Dream Homes Realty",
		ImageUrls:     []string{"https://example.com/images/property1.jpg", "https://example.com/images/property2.jpg"},
		PageDetailUrl: "https://example.com/property/12345",
		Description:   "A beautiful 2-bedroom apartment with city views, close to all amenities.",
		Size:          80,
		Features:      datatypes.JSON([]byte(jsonStringValue)),
		OriginUrl:     "https://example.com/property/12345",
	}

	count := 1
	if n != 0 {
		count = n
	}

	properties := make([]models.Property, count)
	for i := range count {
		newProperty := mergeProperties(defaultProperty, record)
		newProperty.CreatedAt = time.Now()
		newProperty.UpdatedAt = time.Now()
		properties[i] = newProperty
	}

	return properties
}

var Filters = map[string]models.ScraperFilter{
	"casa": {
		Area:           "01f0eef0", // bologna comune
		ContractType:   models.RENT,
		PlatformSource: models.CASA,
		MaxPrice:       1500,
		ToPage:         1,
	},
	"idealista": {
		Area:           "bologna-provincia",
		ContractType:   models.RENT,
		PlatformSource: models.IDEALISTA,
		MaxPrice:       1500,
		FromSize:       80,
		ToPage:         1,
	},
	"immobiliare": {
		Area:           "bologna",
		ContractType:   models.RENT,
		PlatformSource: models.IMMOBILIARE,
		MaxPrice:       1500,
		FromSize:       80,
		ToPage:         1,
	},
}

var Users = map[string]models.User{
	"example": {
		Username: "example@gmail.com",
		Password: "password",
	},
}
