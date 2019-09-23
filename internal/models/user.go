package models

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Name     string
	Email    string     `gorm:"type:varchar(100);unique_index"`
	Accounts []*Account `gorm:"foreignkey:UserID"`
}
