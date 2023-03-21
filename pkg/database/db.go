package db

import (
	"log"
	"os"

	models "github.com/hacktues-9/tf-api/pkg/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func Init() *gorm.DB {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")
	password := os.Getenv("DB_PASSWORD")
	sslmode := os.Getenv("DB_SSLMODE")
	dsn := "host=" + host + " user=" + user + " password=" + password + " dbname=" + dbname + " port=" + port + " sslmode=" + sslmode
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(&models.Projects{}, &models.Creators{}, &models.Pictures{}, &models.Votes{})
	if err != nil {
		log.Fatal(err)
		return
	}
}

func Drop(db *gorm.DB) {
	err := db.Migrator().DropTable(&models.Projects{}, &models.Creators{}, &models.Pictures{}, &models.Votes{})
	if err != nil {
		log.Fatal(err)
		return
	}
}
