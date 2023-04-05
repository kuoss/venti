package store

import (
	"fmt"
	"log"

	"github.com/kuoss/venti/pkg/auth"
	"github.com/kuoss/venti/pkg/configuration"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// todo remove
// const dbfilepath = "./data/venti.sqlite3"

type UserStore struct {
	db *gorm.DB
}

func NewUserStore(filepath string, config configuration.UsersConfig) (*UserStore, error) {
	log.Println("Initializing database...")

	db, err := gorm.Open(sqlite.Open(filepath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("cannot open database: %w", err)
	}

	err = db.AutoMigrate(auth.User{})
	if err != nil {
		return nil, fmt.Errorf("auto migration failed: %w", err)
	}
	setEtcUsers(db, config)
	return &UserStore{db}, nil
}

func setEtcUsers(db *gorm.DB, config configuration.UsersConfig) {

	for _, etcUser := range config.EtcUsers {
		var user auth.User
		result := db.First(&user, "username = ?", etcUser.Username)
		if result.RowsAffected == 0 {
			db.Create(&auth.User{Username: etcUser.Username, Hash: etcUser.Hash, IsAdmin: etcUser.IsAdmin})
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

func (s *UserStore) FindByUsername(name string) (auth.User, error) {
	var user auth.User
	tx := s.db.First(&user, "username = ?", name)
	return user, tx.Error
}

func (s *UserStore) FindByUserIdAndToken(id, token string) (auth.User, error) {
	var user auth.User
	tx := s.db.First(&user, "ID = ? AND token = ?", id, token)
	return user, tx.Error
}

func (s *UserStore) Save(user auth.User) error {
	return s.db.Save(&user).Error
}
