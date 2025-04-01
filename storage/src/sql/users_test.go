package sql

import (
	"testing"

	"github.com/70X/properties-seeker/storage"
	"github.com/70X/properties-seeker/storage/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UsersTestSuite struct {
	suite.Suite
	conn    *storage.DB
	sql     *SQLUsers
	sqlProp *SQLProperties
}

// SetupSuite runs once before all tests in the suite
func (suite *UsersTestSuite) SetupSuite() {
	suite.conn = tests.InitTestDb()
	suite.sql = NewSQLUsers(suite.conn)
	suite.sqlProp = NewSQLProperties(suite.conn)
}

// SetupTest runs before each test in the suite
func (suite *UsersTestSuite) SetupTest() {
	tests.BeforeEachTest(suite.conn)
}

func (suite *UsersTestSuite) TestSQLUsers_CreateUser_Success() {
	user := tests.Users["example"]

	u, _ := suite.sql.CreateUser(user)

	assert.Equal(suite.T(), 1, int(u.ID))
}

func (suite *UsersTestSuite) TestSQLUsers_GetUserByUsername_Success() {
	user := tests.Users["example"]
	suite.sql.CreateUser(user)

	u, _ := suite.sql.GetUserByUsername(user.Username)

	assert.Equal(suite.T(), 1, int(u.ID))
}

func (suite *UsersTestSuite) TestSQLUsers_GetUser_Success() {
	user := tests.Users["example"]
	suite.sql.CreateUser(user)

	u, _ := suite.sql.GetUser(user)

	assert.Equal(suite.T(), 1, int(u.ID))
}

func (suite *UsersTestSuite) TestSQLUsers_GetUserByUsernameNotFound_Success() {
	user := tests.Users["example"]

	u, _ := suite.sql.GetUserByUsername(user.Username)

	assert.Empty(suite.T(), u)
}

func (suite *UsersTestSuite) TestSQLUsers_GetUserNotFound_Success() {
	user := tests.Users["example"]

	u, _ := suite.sql.GetUser(user)

	assert.Empty(suite.T(), u)
}

func TestUsersTestSuite(t *testing.T) {
	suite.Run(t, new(UsersTestSuite))
}
