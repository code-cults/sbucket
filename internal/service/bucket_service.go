package service

import (
	"errors"
	"github.com/srirangamuc/sbucket/internal/db"
	"github.com/srirangamuc/sbucket/internal/model"
)

func CreateBucket(ownerID int,name string) (*model.Bucket, error){
	var existing model.Bucket
	if err := db.DB.Where("name = ?",name).First(&existing).Error; err == nil{
		return nil, errors.New("Bucket already exists")
	}

	bucket := model.Bucket{
		Name : name,
		OwnerID : uint(ownerID),
	}

	if err := db.DB.Create(&bucket).Error ; err != nil{
		return nil, err
	}

	return &bucket, nil
}