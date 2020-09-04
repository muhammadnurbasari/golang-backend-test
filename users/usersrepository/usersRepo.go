package usersrepository

import (
	"golang-backend-test/models/usersmodel"
	"golang-backend-test/users"

	"github.com/jinzhu/gorm"
)

type sqlRepository struct {
	Conn *gorm.DB
}

func NewUsersRepository(Conn *gorm.DB) users.UsersRepository {
	return &sqlRepository{Conn}
}

// RetrievePasswordByName - get password & token to Login method
func (db *sqlRepository) RetrievePasswordByName(userName string) (*usersmodel.ResRetrievePass, error) {
	var users usersmodel.Users

	err := db.Conn.Select("users_id, password").Where("username = ?", userName).Find(&users).Error

	if err != nil {
		return nil, err
	}

	resp := usersmodel.ResRetrievePass{
		UsersId:  users.UsersId,
		Password: users.Password,
	}

	return &resp, nil
}
