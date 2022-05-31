package datamodel

import "gorm.io/gorm"

type Todo struct {
	gorm.Model
	Text   string
	Done   bool
	UserID uint
}
