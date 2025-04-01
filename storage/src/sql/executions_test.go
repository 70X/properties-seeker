package sql

import (
	"testing"

	"github.com/70X/properties-seeker/storage"
	"github.com/70X/properties-seeker/storage/src/models"
	"github.com/70X/properties-seeker/storage/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ExecutionsTestSuite struct {
	suite.Suite
	conn    *storage.DB
	sql     *SQLExecutions
	sqlProp *SQLProperties
}

func (suite *ExecutionsTestSuite) SetupSuite() {
	suite.conn = tests.InitTestDb()
	suite.sql = NewSQLExecutions(suite.conn)
	suite.sqlProp = NewSQLProperties(suite.conn)
}

func (suite *ExecutionsTestSuite) SetupTest() {
	tests.BeforeEachTest(suite.conn)
}

func (suite *ExecutionsTestSuite) TestSQLExecutions_AddExecution_Success() {
	scraperFilter := tests.Filters["immobiliare"]

	e, _ := suite.sql.StartExecution(scraperFilter)

	assert.Equal(suite.T(), 1, int(e.ID))
}

func (suite *ExecutionsTestSuite) TestSQLExecutions_GetExecutionByFilter_Success() {
	scraperFilter := models.ScraperFilter{
		ContractType:   "SELL",
		Area:           "Downtown",
		PlatformSource: "IMMOBILIARE",
		FromSize:       50,
		MaxPrice:       500000,
		FromPage:       1,
		ToPage:         10,
	}
	suite.sql.StartExecution(scraperFilter)

	execution, _ := suite.sql.GetExecutionByFilter(scraperFilter)
	assert.Equal(suite.T(), 1, int(execution.ID))
	assert.Equal(suite.T(), models.PROGRESS, execution.Status)
	assert.Equal(suite.T(), "Downtown", execution.Area)
	assert.Equal(suite.T(), models.IMMOBILIARE, execution.PlatformSource)
	assert.Equal(suite.T(), models.SELL, execution.ContractType)
	assert.Equal(suite.T(), 0, execution.Results)
}

func (suite *ExecutionsTestSuite) TestSQLExecutions_SetExecutionByFilter_Success() {
	scraperFilter := models.ScraperFilter{
		ContractType:   "SELL",
		Area:           "Downtown",
		PlatformSource: "IMMOBILIARE",
		FromSize:       50,
		MaxPrice:       500000,
		FromPage:       1,
		ToPage:         10,
	}
	scraperFilter2 := scraperFilter
	scraperFilter2.ContractType = models.RENT
	e, _ := suite.sql.StartExecution(scraperFilter)
	suite.sql.StartExecution(scraperFilter2)

	suite.sql.SetExecutionStatus(e.ID, models.SUCCESS, 25, 25)

	execution, _ := suite.sql.GetExecutionByFilter(scraperFilter)
	assert.Equal(suite.T(), models.SUCCESS, execution.Status)
}

func (suite *ExecutionsTestSuite) TestSQLExecutions_GetExecutionByFilter_NotFound() {
	scraperFilter := tests.Filters["immobiliare"]

	execution, err := suite.sql.GetExecutionByFilter(scraperFilter)

	assert.Equal(suite.T(), "record not found", err.Error())
	assert.Equal(suite.T(), 0, int(execution.ID))
}

// func (suite *ExecutionsTestSuite) TestSQLExecutions_AddExecutionAssignedToProperty_Success() {
// 	scraperFilter := models.ScraperFilter{
// 		ContractType:   "SELL",
// 		Area:           "Downtown",
// 		PlatformSource: "IMMOBILIARE",
// 		FromSize:       50,
// 		MaxPrice:       500000,
// 		FromPage:       1,
// 		ToPage:         10,
// 	}
// 	execution, _ := suite.sql.StartExecution(scraperFilter)
// 	property := models.Property{
// 		Title:      "Cozy 2-Bedroom Apartment in Downtown",
// 		PlatformID: "P123456",
// 	}
// 	propIds, _ := suite.sqlProp.FetchSellProperties([]common.Property{property}, execution.ID)

// 	p, _ := suite.sqlProp.GetProperty(uint(propIds[0].ID))

// 	assert.Equal(suite.T(), execution.ID, p.Filter.ID)
// }

func TestExecutionsTestSuite(t *testing.T) {
	suite.Run(t, new(ExecutionsTestSuite))
}
