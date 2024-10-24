package models

import "gorm.io/gorm"

type Contact struct {
	gorm.Model
	Id       int    `json:"id" gorm:"primary_key"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}
