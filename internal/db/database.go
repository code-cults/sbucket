package db
 
import (
	"log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/srirangamuc/sbucket/internal/config"
)

var DB *gorm.DB

func Connect() {
	dsn := config.GetEnv("DB_URL", "")
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
	}
	DB = database
	log.Println("Connected to Database successfully")
}