package usersmodel

type Users struct {
	// gorm.Model
	UsersId  int
	Username string
	Password string
	Role     int
}

type ResRetrievePass struct {
	UsersId  int    `json:"users_id"`
	Password string `json:"password"`
}

type ReqLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ResLogin struct {
	UsersId int    `json:"users_id"`
	Token   string `json:"token"`
}
