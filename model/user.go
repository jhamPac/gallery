package model

import (
	"errors"

	"github.com/jhampac/gallery/hasho"
	"github.com/jhampac/gallery/rando"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrNotFound        = errors.New("model: resource not found")
	ErrInvalidID       = errors.New("model: ID provided was invalid")
	ErrInvalidPassword = errors.New("model: incorrect password provided")
)

const (
	pepper        = "suns-in-7"
	hmacSecretKey = "change-this-secert-later-for-production"
)

type UserService struct {
	UserDB
}

func NewUserService(connInfo string) (*UserService, error) {
	ug, err := newUserGorm(connInfo)
	if err != nil {
		return nil, err
	}

	return &UserService{
		UserDB: ug,
	}, nil
}

type UserDB interface {
	// db look ups with args
	ByID(id uint) (*User, error)
	ByEmail(email string) (*User, error)
	ByRemember(token string) (*User, error)

	// altering db data
	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error

	// db operations
	Close() error
	AutoMigrate() error
	DestructiveReset() error
}

// user db representation
type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;unique_index"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
	Remember     string `gorm:"-"`
	RememberHash string `gorm:"not null;unique_index"`
}

type userGorm struct {
	db   *gorm.DB
	hmac hasho.HMAC
}

func newUserGorm(connInfo string) (*userGorm, error) {
	db, err := gorm.Open("postgres", connInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)

	hmac := hasho.NewHMAC(hmacSecretKey)

	return &userGorm{
		db:   db,
		hmac: hmac,
	}, nil
}

func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}

func (ug *userGorm) Create(user *User) error {
	hashedBytes, err := bcrypt.GenerateFromPassword(
		[]byte(user.Password+pepper), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// hash password
	user.PasswordHash = string(hashedBytes)
	user.Password = ""

	// check remember token
	if user.Remember == "" {
		token, err := rando.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
	}
	user.RememberHash = ug.hmac.Hash(user.Remember)

	return ug.db.Create(user).Error
}

func (ug *userGorm) ByID(id uint) (*User, error) {
	var user User
	db := ug.db.Where("id = ?", id)
	err := first(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ug *userGorm) ByEmail(email string) (*User, error) {
	var user User
	db := ug.db.Where("email = ?", email)
	err := first(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ug *userGorm) ByRemember(token string) (*User, error) {
	var user User
	rememberHash := ug.hmac.Hash(token)
	err := first(ug.db.Where("remember_hash = ?", rememberHash), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ug *userGorm) Update(user *User) error {
	if user.Remember != "" {
		user.RememberHash = ug.hmac.Hash(user.Remember)
	}
	return ug.db.Save(user).Error
}

func (ug *userGorm) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}
	user := User{Model: gorm.Model{ID: id}}
	return ug.db.Delete(&user).Error
}

func (ug *userGorm) Close() error {
	return ug.db.Close()
}

func (ug *userGorm) DestructiveReset() error {
	err := ug.db.DropTableIfExists(&User{}).Error
	if err != nil {
		return err
	}
	return ug.AutoMigrate()
}

func (ug *userGorm) AutoMigrate() error {
	if err := ug.db.AutoMigrate(&User{}).Error; err != nil {
		return err
	}
	return nil
}

func (us *UserService) Authenticate(email string, password string) (*User, error) {
	foundUser, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(password+pepper))
	switch err {
	case nil:
		return foundUser, nil
	case bcrypt.ErrMismatchedHashAndPassword:
		return nil, ErrInvalidPassword
	default:
		return nil, err
	}
}
