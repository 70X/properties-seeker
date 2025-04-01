package sql

import (
	"errors"
	"testing"
	"time"

	"github.com/70X/properties-seeker/storage"
	"github.com/70X/properties-seeker/storage/src/models"
	"github.com/70X/properties-seeker/storage/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type PropertiesTestSuite struct {
	suite.Suite
	conn          *storage.DB
	sql           *SQLProperties
	sqlExecutions *SQLExecutions
	executionId   uint
}

// SetupSuite runs once before all tests in the suite
func (suite *PropertiesTestSuite) SetupSuite() {
	suite.conn = tests.InitTestDb()
	suite.sql = NewSQLProperties(suite.conn)
	suite.sqlExecutions = NewSQLExecutions(suite.conn)
}

// SetupTest runs before each test in the suite
func (suite *PropertiesTestSuite) SetupTest() {
	tests.BeforeEachTest(suite.conn)
	e, _ := suite.sqlExecutions.StartExecution(tests.Filters["immobiliare"])
	suite.executionId = e.ID
}

func (suite *PropertiesTestSuite) TestSQLProperties_AddProperties_Success() {
	properties := tests.NewNProperties(models.Property{}, 0)

	res, _ := suite.sql.AddProperties(properties)

	assert.Equal(suite.T(), 1, len(res))
	assert.Equal(suite.T(), 1, int(res[0].ID))
}

func (suite *PropertiesTestSuite) TestSQLProperties_Add5Properties_Success() {
	properties := tests.NewNProperties(models.Property{}, 5)

	res, _ := suite.sql.AddProperties(properties)

	assert.Equal(suite.T(), 5, len(res))
	assert.Equal(suite.T(), 1, int(res[0].ID))
	assert.Equal(suite.T(), 5, int(res[4].ID))
}

func (suite *PropertiesTestSuite) TestSQLProperties_NotAddDuplicatedProperty_Success() {
	properties := tests.NewNProperties(models.Property{}, 3)
	properties[0].PlatformID = "P2222222"
	properties[1].PlatformID = "P1111111"
	properties[2].PlatformID = "P2222222"

	res, _ := suite.sql.AddProperties(properties)

	assert.Equal(suite.T(), 2, len(res))
	assert.Equal(suite.T(), "P2222222", res[0].PlatformID)
	assert.Equal(suite.T(), "P1111111", res[1].PlatformID)
}

func (suite *PropertiesTestSuite) TestSQLProperties_GetProperty_NotFound() {
	res, err := suite.sql.GetProperty(uint(1))

	assert.Equal(suite.T(), errors.New("record not found"), err)
	assert.Nil(suite.T(), res)
}

func (suite *PropertiesTestSuite) TestSQLProperties_GetProperty_Success() {
	properties := tests.NewNProperties(models.Property{}, 3)
	res, _ := suite.sql.AddProperties(properties)

	property, err := suite.sql.GetProperty(res[0].ID)

	assert.Equal(suite.T(), res[0].PlatformID, property.PlatformID)
	assert.Nil(suite.T(), err)
}

func (suite *PropertiesTestSuite) TestSQLProperties_GetPropertyWithExecution_Success() {
	e, _ := suite.sqlExecutions.StartExecution(tests.Filters["idealista"])
	properties := tests.NewNProperties(models.Property{ExecutionID: e.ID}, 1)
	res, _ := suite.sql.AddProperties(properties)

	property, _ := suite.sql.GetProperty(res[0].ID)

	assert.Equal(suite.T(), property.Execution.ScraperFilter.Area, tests.Filters["idealista"].Area)
	assert.Equal(suite.T(), property.Execution.ToPage, tests.Filters["idealista"].ToPage)
	assert.Equal(suite.T(), property.Execution.PlatformSource, tests.Filters["idealista"].PlatformSource)
}

func (suite *PropertiesTestSuite) TestSQLProperties_FetchPropertiesWithArchiveFlag_Success() {
	properties := tests.NewNProperties(models.Property{}, 1)
	property := properties[0]
	t := time.Now()
	property.ArchivedAt = &t

	_, err := suite.sql.UpdateProperty(property)

	assert.Nil(suite.T(), err)
}

func (suite *PropertiesTestSuite) TestSQLProperties_FetchRentProperties_Success() {
	filter := tests.Filters["idealista"]
	filter.ContractType = models.RENT
	e, _ := suite.sqlExecutions.StartExecution(filter)
	properties := tests.NewNProperties(models.Property{ExecutionID: e.ID}, 10)
	suite.sql.AddProperties(properties)

	list, _ := suite.sql.FetchRentProperties(false)

	assert.Equal(suite.T(), 10, len(list))
}

func (suite *PropertiesTestSuite) TestSQLProperties_FetchSellProperties_Success() {
	filter := tests.Filters["idealista"]
	filter.ContractType = models.SELL
	e, _ := suite.sqlExecutions.StartExecution(filter)
	properties := tests.NewNProperties(models.Property{ExecutionID: e.ID}, 10)
	suite.sql.AddProperties(properties)

	list, _ := suite.sql.FetchSellProperties(false)

	assert.Equal(suite.T(), 10, len(list))
}

func (suite *PropertiesTestSuite) TestSQLProperties_SetPropertyNote_Success() {
	properties := tests.NewNProperties(models.Property{}, 1)
	ids, _ := suite.sql.AddProperties(properties)
	properties[0].Note = "A test note"

	suite.sql.UpdateProperty(properties[0])

	item, _ := suite.sql.GetProperty(ids[0].ID)
	assert.Equal(suite.T(), "A test note", item.Note)
}

func (suite *PropertiesTestSuite) TestSQLProperties_FetchPropertiesWithAuction_Success() {
	properties := tests.NewNProperties(models.Property{IsAuction: true}, 1)
	suite.sql.AddProperties(properties)

	res, _ := suite.sql.FetchAuctionList(false)

	assert.Equal(suite.T(), 1, len(res))
	assert.Equal(suite.T(), 1, int(res[0].ID))
}

func (suite *PropertiesTestSuite) TestSQLProperties_FetchPropertiesWithAuctionArchived_Success() {
	properties := tests.NewNProperties(models.Property{IsAuction: true}, 1)
	suite.sql.AddProperties(properties)
	t := time.Now()
	properties[0].ArchivedAt = &t
	suite.sql.UpdateProperty(properties[0])

	res, _ := suite.sql.FetchAuctionList(true)

	assert.Equal(suite.T(), 1, len(res))
	assert.Equal(suite.T(), 1, int(res[0].ID))
}

func (suite *PropertiesTestSuite) TestSQLProperties_DeleteProperty_Success() {
	properties := tests.NewNProperties(models.Property{}, 1)
	suite.sql.AddProperties(properties)

	suite.sql.DeleteProperty(properties[0].ID)

	_, err := suite.sql.GetProperty(properties[0].ID)
	assert.Equal(suite.T(), "record not found", err.Error())
}

func (suite *PropertiesTestSuite) TestSQLProperties_GetFilterProperty_Success() {
	properties := tests.NewNProperties(models.Property{}, 1)
	suite.sql.AddProperties(properties)

	list, _ := suite.sql.FetchRentProperties(false)

	assert.Equal(suite.T(), tests.Filters["immobiliare"].Area, list[0].Execution.Area)
}

func TestPropertiesTestSuite(t *testing.T) {
	suite.Run(t, new(PropertiesTestSuite))
}
