package model

import (
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type User struct {
	Id         uint   `gorm:"primarykey;autoIncrement:true"`
	PassSerie  uint32 `gorm:"uniqueIndex:pass"`
	PassNumber uint32 `gorm:"uniqueIndex:pass"`
	Name       string
	Surname    string
	Patronimic string
	Address    string
	Tasks      []Task `gorm:"foreignKey:UserId"`
}

type Task struct {
	Id        uint `gorm:"primarykey;autoIncrement:true"`
	Desc      string
	TimeStart *time.Time
	TimeEnd   *time.Time
	UserId    uint
}

var db *gorm.DB

func init() {
	link, ok := os.LookupEnv("DATABASE_LINK")
	if !ok {
		panic("DATABASE_LINK environment variable is undefined")
	}
	var err error
	db, err = gorm.Open(postgres.Open(link), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("Failed to connect to database")
	}
	db.AutoMigrate(&User{}, &Task{})
}

func GetDB() *gorm.DB {
	return db
}
