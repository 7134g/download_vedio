package dao

import "gorm.io/gorm"

var db *gorm.DB

func DaoInit(db *gorm.DB) {
	db = db
}
