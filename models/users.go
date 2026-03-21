package models

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	FullName  string         `gorm:"size:255;not null" json:"full_name"`
	Username  string         `gorm:"size:100;uniqueIndex;not null" json:"username"`
	Email     string         `gorm:"size:255;uniqueIndex;not null" json:"email"`
	AvatarURL *string        `gorm:"size:255" json:"avatar_url,omitempty"`
	Password  string         `gorm:"not null" json:"-"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}

	return
}

func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	if u.Password != "" && len(u.Password) < 60 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)

		if err != nil {
			return err
		}

		u.Password = string(hashedPassword)
	}

	return
}
