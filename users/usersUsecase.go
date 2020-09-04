package users

import "golang-backend-test/models/usersmodel"

type UsersUsecase interface {
	LoginUsers(dataRequest *usersmodel.ReqLogin) (*usersmodel.ResLogin, error)
}
