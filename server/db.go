package server

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB() {
	log.Println("Initializing database...")

	// open db
	var err error
	db, err = gorm.Open(sqlite.Open("./data/venti.sqlite3"), &gorm.Config{})
	if err != nil {
		log.Fatal("Cannot open DB.")
	}

	// Database Migrations
	db.AutoMigrate(&User{})

	// Database Seeding
	// etc users
	for _, etcUser := range GetConfig().EtcUsersConfig.EtcUsers {
		var user User
		result := db.First(&user, "username = ?", etcUser.Username)
		if result.RowsAffected == 0 {
			db.Create(&User{Username: etcUser.Username, Hash: etcUser.Hash, IsAdmin: etcUser.IsAdmin})
			log.Println("User '" + etcUser.Username + "' added.")
		} else {
			log.Println("User '" + etcUser.Username + "' already exists.")
			if user.Hash != etcUser.Hash {
				user.Hash = etcUser.Hash
				db.Save(&user)
				log.Println("User '" + etcUser.Username + "' updated.")
			}
		}
	}
}
