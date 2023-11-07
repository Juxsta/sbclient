package main

import (
	"log"

	qbtp "github.com/Juxsta/sbclient/seedbox-service/qbittorrentproxy"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Username  string
	Password  string
	SessionID string
}

type Category struct {
	gorm.Model
	qbtp.Category
}

func initDB() *gorm.DB {
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
			user = &User{Model: gorm.Model{ID: 1}, Username: "admin", Password: string(hash), SessionID: "tjZcBehXWvgh4Eb6ilHmlhkFEsd2nGfu"}
			db.Create(user)
		} else {
			log.Fatalf("Error retrieving user: %v", err)
		}
	}

	return db
}

func getCategories(m *MyServer) map[string]qbtp.Category {
	var categories []qbtp.Category
	m.db.Find(&categories)

	catMap := make(map[string]qbtp.Category)
	for _, category := range categories {
		catMap[category.Category] = category
	}
	return catMap
}
