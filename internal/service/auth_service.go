package service

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"github.com/srirangamuc/sbucket/internal/db"
	"github.com/srirangamuc/sbucket/internal/model"
	"time"
	"os"
	"github.com/golang-jwt/jwt/v5"
)

func SignUpUser(email,password string) error {
	var existing model.User
	result := db.DB.Where("email=?",email).First(&existing)
	if result.Error == nil {
		return errors.New("User Already Exists")
	}

	hash,err := bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost)
	if err!=nil{
		return err
	}

	user := model.User{
		Email : email,
		PasswordHash : string(hash),
	}

	if err:= db.DB.Create(&user).Error; err != nil{
		return err
	}

	return nil
}

func LoginUser(email,password string) (string,error){
	var user model.User
	result := db.DB.Where("email = ?",email).First(&user)
	if result.Error != nil{
		return "", errors.New("invalid credentials")
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash),[]byte(password))
	if err != nil{
		return "",errors.New("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour*72).Unix(),
	})

	secret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil{
		return "",err
	}
	return tokenString,nil
}