package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func SetupModels() *gorm.DB {
	db, err := gorm.Open("mysql", "root:.kokoko1@(localhost:3306)/golangdb?charset=utf8&parseTime=True&loc=Local")

	if err != nil {
		panic("gagal koneksi database")
	}

	db.AutoMigrate(&Nasabah{})

	return db
}
