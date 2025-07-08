package db

import (
	"fmt"
	"log"
	"medods_testcase/db/models"
	"medods_testcase/utils"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Prepare(db *gorm.DB) {
	err := db.AutoMigrate(&models.User{}, &models.Token{})

	if err != nil {
		log.Panic("Ошибка автоматической миграции")
	}

	log.Println("Выполнена автоматическая миграция")
}

func NewDbCon() *gorm.DB {
	config := utils.NewConfig()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", config.DbHost, config.DbUser, config.DbPass, config.DbName, config.DbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalln("Ошибка подключения к базе данных", err)
	}

	log.Println("Соединение с базой данных установлено")

	Prepare(db)

	return db
}
