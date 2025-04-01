package sql

import (
	"testing"

	"github.com/70X/properties-seeker/storage"
	"github.com/70X/properties-seeker/storage/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ScraperFilterTestSuite struct {
	suite.Suite
	conn *storage.DB
	sql  *SQLScraperFilter
}

func (suite *ScraperFilterTestSuite) SetupSuite() {
	suite.conn = tests.InitTestDb()
	suite.sql = NewSQLScraperFilter(suite.conn)
}

func (suite *ScraperFilterTestSuite) SetupTest() {
	tests.BeforeEachTest(suite.conn)
}

func (suite *ScraperFilterTestSuite) TestSQLScraperFilters_AddScraperFilter_Success() {
	scraperFilter := tests.Filters["immobiliare"]

	id, _ := suite.sql.AddScraperFilter(scraperFilter)

	assert.Equal(suite.T(), 1, int(id))
}

func (suite *ScraperFilterTestSuite) TestSQLScraperFilters_AddScraperFilterEnabledByDefault_Success() {
	scraperFilter := tests.Filters["immobiliare"]
	suite.sql.AddScraperFilter(scraperFilter)

	filters, _ := suite.sql.GetScraperFilters()

	assert.NotEmpty(suite.T(), filters[0])
}

func (suite *ScraperFilterTestSuite) TestSQLScraperFilters_SetScraperFilterDisabled_Success() {
	scraperFilter := tests.Filters["immobiliare"]
	id, _ := suite.sql.AddScraperFilter(scraperFilter)

	_, err := suite.sql.EnableScraperFilter(id, false)

	assert.Empty(suite.T(), err)
}

func (suite *ScraperFilterTestSuite) TestSQLScraperFilters_GetEmptyScraperFilterIfDisabled_Success() {
	scraperFilter := tests.Filters["immobiliare"]
	id, _ := suite.sql.AddScraperFilter(scraperFilter)
	suite.sql.EnableScraperFilter(id, false)

	filters, _ := suite.sql.GetScraperFilters()

	assert.Empty(suite.T(), filters)
}

func (suite *ScraperFilterTestSuite) TestSQLScraperFilters_GetScraperFilterRightType_Success() {
	scraperFilter := tests.Filters["immobiliare"]
	suite.sql.AddScraperFilter(scraperFilter)

	filters, _ := suite.sql.GetScraperFilters()

	assert.Equal(suite.T(), scraperFilter.Area, filters[0].Area)
	assert.Equal(suite.T(), scraperFilter.ContractType, filters[0].ContractType)
}

func TestScraperFilterTestSuite(t *testing.T) {
	suite.Run(t, new(ScraperFilterTestSuite))
}
