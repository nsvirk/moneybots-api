package api

import (
	"crypto/md5"
	"encoding/hex"
	"time"
)

// Users struct represents the user model
type UserModel struct {
	ID           uint   `gorm:"primarykey"`
	UserId       string `gorm:"type:varchar(10);index;unique;not null"`
	PasswordHash string `gorm:"type:varchar(100)"`
	Enctoken     string
	LoginTime    string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// TableName returns the table name of the user model
func (u *UserModel) TableName() string {
	return "api_users"
}

// generateMD5Hash generates an MD5 hash from a string
func generateMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
