package database

import (
	"choccy/server/database/model"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"os"
)

var DB *gorm.DB

func InitDB(path string) bool {
	isNew := false
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		isNew = true
	} else if err != nil {
		panic(err.Error())
	}

	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic("failed to connect database")
	}
	DB = db
	err = DB.AutoMigrate(&model.Project{}, &model.Setting{}, &model.Task{}, &model.TaskResult{}, &model.Status{})
	if err != nil {
		panic(err.Error())
	}

	return isNew
}

func Paginate(page int, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}
		if pageSize <= 0 {
			pageSize = 20
		}
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
