package user

import (
	"fmt"
	"log"

	"github.com/kuoss/common/logger"
	"github.com/kuoss/venti/pkg/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// todo remove
// const dbfilepath = "./data/venti.sqlite3"

type UserService struct {
	db *gorm.DB
}

func New(filepath string, config model.UserConfig) (*UserService, error) {
	log.Println("Initializing database...")

	db, err := gorm.Open(sqlite.Open(filepath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("DB open err: %w", err)
	}

	err = db.AutoMigrate(model.User{})
	if err != nil {
		return nil, fmt.Errorf("auto migration failed: %w", err)
	}
	setEtcUsers(db, config)
	return &UserService{db}, nil
}

func setEtcUsers(db *gorm.DB, config model.UserConfig) {

	for _, etcUser := range config.EtcUsers {
		var user model.User
		result := db.First(&user, "username = ?", etcUser.Username)
		if result.RowsAffected == 0 {
			db.Create(&model.User{Username: etcUser.Username, Hash: etcUser.Hash, IsAdmin: etcUser.IsAdmin})
			logger.Infof("User '%s' added.", etcUser.Username)
		} else {
			logger.Infof("User '%s' already exists.", etcUser.Username)
			if user.Hash != etcUser.Hash {
				user.Hash = etcUser.Hash
				db.Save(&user)
				logger.Infof("User '%s' updated.", etcUser.Username)
			}
		}
	}
}

func (s *UserService) FindByUsername(name string) (model.User, error) {
	var user model.User
	tx := s.db.First(&user, "username = ?", name)
	return user, tx.Error
}

func (s *UserService) FindByUserIdAndToken(id, token string) (model.User, error) {
	var user model.User
	tx := s.db.First(&user, "ID = ? AND token = ?", id, token)
	return user, tx.Error
}

func (s *UserService) Save(user model.User) error {
	return s.db.Save(&user).Error
}
