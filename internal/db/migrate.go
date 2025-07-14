package db

import (
	"log"
	"github.com/srirangamuc/sbucket/internal/model"
)

func AutoMigrateModels(){
	err := DB.AutoMigrate(
		&model.User{},
		&model.Bucket{},
	)

	if err != nil{
		log.Fatal("Failed to migrate: %v",err)
	}

	log.Println("Database migrated successfully")
}