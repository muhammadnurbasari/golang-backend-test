package usersusecase

import (
	"errors"
	"golang-backend-test/jwt"
	"golang-backend-test/models/usersmodel"
	"golang-backend-test/users"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UsersUsecase struct {
	usersRepository users.UsersRepository
}

func NewUsersUsecase(usersRepo users.UsersRepository) users.UsersUsecase {
	return &UsersUsecase{
		usersRepository: usersRepo,
	}
}

func (usersUC *UsersUsecase) LoginUsers(dataRequest *usersmodel.ReqLogin) (*usersmodel.ResLogin, error) {
	retrievePassword, err := usersUC.usersRepository.RetrievePasswordByName(dataRequest.Username)

	if err != nil {
		return nil, errors.New("1")
	}

	if retrievePassword == nil {
		return nil, errors.New("1")
	}

	errCheckPass := bcrypt.CompareHashAndPassword([]byte(retrievePassword.Password), []byte(dataRequest.Password))

	if errCheckPass != nil {
		return nil, errors.New("2")
	}

	timeExpDuration := time.Duration(60*1) * time.Minute
	token, errToken := jwt.GenerateTokenJwt(int(retrievePassword.UsersId), timeExpDuration)

	if errToken != nil {
		return nil, errToken
	}

	resp := usersmodel.ResLogin{
		UsersId: retrievePassword.UsersId,
		Token:   token,
	}

	return &resp, nil
}
