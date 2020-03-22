package models

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"html"
	"log"
	"time"
	"strings"
	"github.com/badoux/checkmail"
)

type User struct {
	*gorm.Model
	FullName string `json:"full_name"`
	Username string `gorm:"size:55;not null;unique" json:"username"`
	Password string `gorm:"size:255;not null" json:"password"`
	Email    string `gorm:"size:100;not null;unique" json:"email"`
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (u *User) BeforeSave() error {
	hashedPassword, err := Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) Prepare() {
	u.ID = 0
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	u.FullName = html.EscapeString(strings.TrimSpace(u.FullName))
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))
	u.Password = html.EscapeString(strings.TrimSpace(u.Password))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
}

func Validate(u *User) error {
	if u.FullName == "" {
		return errors.New("Required Full Name")
	}
	if u.Username == "" {
		return errors.New("Required Username")
	}
	if u.Password == "" {
		return errors.New("Required Password")
	}
	if u.Email == "" {
		return errors.New("Required Email")
	}
	err := checkmail.ValidateFormat(u.Email)
	if err != nil {
		return errors.New("Invalid Email")
	}
	return nil
}

func (u *User) ValidateUser(action string) error {
	switch strings.ToLower(action) {
	case "update":
		return Validate(u)
	case "login":
		if u.Email == "" && u.Username == "" {
			return errors.New("Required Email or Username")
		}
		if u.Email != "" {
			err := checkmail.ValidateFormat(u.Email)
			if err != nil {
				return errors.New("Invalid Email")
			}
		}
		if u.Password == "" {
			return errors.New("Required Password")
		}
	default:
		return Validate(u)
	}
	return nil
}

func (u *User) SaveUser(db *gorm.DB) (*User, error) {
	err := db.Debug().Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) FindAllUser(db *gorm.DB) (*[]User, error) {
	var users []User
	err := db.Debug().Model(&User{}).Limit(100).Find(&users).Error
	if err != nil {
		return &[]User{}, err
	}
	return &users, nil
}

func (u *User) FindUserById(db *gorm.DB, uid uint32) (*User, error) {
	err := db.Debug().Model(&User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &User{}, errors.New("User Not Found")
	}
	return u, nil
}

func (u *User) UpdateUserById(db *gorm.DB, uid uint32) (*User, error) {
	err := u.BeforeSave()
	if err != nil {
		log.Fatal(err)
	}
	err = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).UpdateColumn(
		map[string]interface{}{
			"username":   u.Username,
			"password":   u.Password,
			"full_name":  u.FullName,
			"email":      u.Email,
			"updated_at": time.Now(),
		}).Error
	if db.Error != nil {
		return &User{}, db.Error
	}
	// This is the display the updated user
	err = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) DeleteUserById(db *gorm.DB, uid uint32) (int64, error) {
	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).Delete(&User{})
	return db.RowsAffected, nil
}
