package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type UserService struct {
	db *gorm.DB
}

type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null;unique_index"`
}

func NewUserService(connectionInfo string) (*UserService, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return &UserService{
		db: db,
	}, nil
}

func (us *UserService) Close() error {
	return us.db.Close()
}
