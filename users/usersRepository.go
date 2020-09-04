package users

import "golang-backend-test/models/usersmodel"

type UsersRepository interface {
	RetrievePasswordByName(userName string) (*usersmodel.ResRetrievePass, error)
}
