package data

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func InitDb(host, usr, pass, dbName string) {
	if host == "" || usr == "" || pass == "" || dbName == "" {
		return
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Asia/Shanghai", host, usr, pass, dbName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	log.Println("数据库链接成功")
	DB = db
	if err != nil {
		log.Panicln(err)
	}
	err = DB.AutoMigrate(&TtsLog{})
	if err != nil {
		log.Panicln(err)
	}
}

type TtsLog struct {
	Text  string
	Audio []byte
}

func Save(text string, content []byte) {
	if DB == nil {
		return
	}
	tts := TtsLog{text, content}
	err := DB.Create(&tts).Error
	if err != nil {
		log.Println(err.Error())
	}
}
