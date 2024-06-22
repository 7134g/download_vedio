package db

import (
	"dv/internel/serve/api/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	db *gorm.DB
)

func InitSqlite(fp string) *gorm.DB {
	var err error
	db, err = gorm.Open(sqlite.Open(fp), &gorm.Config{
		PrepareStmt: true,
		Logger:      logger.Default.LogMode(logger.Silent),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}

	_ = db.AutoMigrate(&model.Task{})
	return db
}

func GetDB() *gorm.DB {
	return db
}
