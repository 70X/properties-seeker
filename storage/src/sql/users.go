package sql

import (
	"github.com/70X/properties-seeker/storage"
	"github.com/70X/properties-seeker/storage/src/models"
	"golang.org/x/crypto/bcrypt"
)

type SQLUsers struct {
	conn *storage.DB
}

func (s SQLUsers) GetUser(user models.User) (*models.User, error) {
	u, err := s.GetUserByUsername(user.Username)
	if err != nil || u == nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword(
		[]byte(u.Password),
		[]byte(user.Password),
	)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s SQLUsers) GetUserByUsername(username string) (*models.User, error) {
	item := models.User{Username: username}
	res := s.conn.Find(&item)
	if item.ID == 0 {
		return nil, res.Error
	}
	return &item, res.Error
}

func (s SQLUsers) CreateUser(user models.User) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(user.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, err
	}
	item := models.User{Username: user.Username, Password: string(hashedPassword)}
	res := s.conn.Create(&item)
	return &item, res.Error
}

func NewSQLUsers(conn *storage.DB) *SQLUsers {
	return &SQLUsers{conn: conn}
}
