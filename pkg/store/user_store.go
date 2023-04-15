package store

import (
	"fmt"
	"log"

	"github.com/glebarez/sqlite"
	"github.com/kuoss/venti/pkg/model"
	"gorm.io/gorm"
)

// todo remove
// const dbfilepath = "./data/venti.sqlite3"

type UserStore struct {
	db *gorm.DB
}

func NewUserStore(filepath string, config model.UsersConfig) (*UserStore, error) {
	log.Println("Initializing database...")

	db, err := gorm.Open(sqlite.Open(filepath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("cannot open database: %w", err)
	}

	err = db.AutoMigrate(model.User{})
	if err != nil {
		return nil, fmt.Errorf("auto migration failed: %w", err)
	}
	setEtcUsers(db, config)
	return &UserStore{db}, nil
}

func setEtcUsers(db *gorm.DB, config model.UsersConfig) {

	for _, etcUser := range config.EtcUsers {
		var user model.User
		result := db.First(&user, "username = ?", etcUser.Username)
		if result.RowsAffected == 0 {
			db.Create(&model.User{Username: etcUser.Username, Hash: etcUser.Hash, IsAdmin: etcUser.IsAdmin})
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

func (s *UserStore) FindByUsername(name string) (model.User, error) {
	var user model.User
	tx := s.db.First(&user, "username = ?", name)
	return user, tx.Error
}

func (s *UserStore) FindByUserIdAndToken(id, token string) (model.User, error) {
	var user model.User
	tx := s.db.First(&user, "ID = ? AND token = ?", id, token)
	return user, tx.Error
}

func (s *UserStore) Save(user model.User) error {
	return s.db.Save(&user).Error
}
