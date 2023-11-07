package database

import (
	"log"

	qbtp "github.com/Juxsta/sbclient/seedbox-service/qbittorrentproxy"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Username string
	Password string
}

type Category struct {
	gorm.Model
	qbtp.Category
}

func InitDB() *gorm.DB {
	db, err := gorm.Open("sqlite3", "database.db")
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	db.AutoMigrate(&User{})
	db.AutoMigrate(&Category{})

	// Set a predefined user and session ID if there is no user with id 1
	user := &User{}
	if err := db.First(user).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			hash, err := bcrypt.GenerateFromPassword([]byte("adminadmin"), bcrypt.DefaultCost)
			if err != nil {
				log.Fatalf("Error generating password hash: %v", err)
			}
			user = &User{Model: gorm.Model{ID: 1}, Username: "admin", Password: string(hash)}
			db.Create(user)
		} else {
			log.Fatalf("Error retrieving user: %v", err)
		}
	}

	return db
}
