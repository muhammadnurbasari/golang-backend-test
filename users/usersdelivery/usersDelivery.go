package usersdelivery

import (
	"golang-backend-test/helpers"
	"golang-backend-test/middleware"
	"golang-backend-test/models/usersmodel"
	"golang-backend-test/users"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type usersHandler struct {
	UsersUsecase users.UsersUsecase
}

func NewUsersHttpHandler(r *gin.Engine, usersUC users.UsersUsecase) {
	handler := usersHandler{
		UsersUsecase: usersUC,
	}

	auth := r.Group("/users")
	auth.POST("/login", middleware.BasicAuth, handler.PostLogin)
}

func (handler *usersHandler) PostLogin(c *gin.Context) {
	var req usersmodel.ReqLogin
	errBind := c.BindJSON(&req)

	if errBind != nil {
		log.Error().Msg(errBind.Error())
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"message": helpers.MsgBadReq,
				"status":  http.StatusBadRequest,
			},
		)
		return
	}

	result, err := handler.UsersUsecase.LoginUsers(&req)

	if err != nil {
		log.Error().Msg(err.Error())

		if err.Error() == "1" {
			c.JSON(
				http.StatusForbidden,
				gin.H{
					"message": "Username is not found on set rows",
					"status":  http.StatusForbidden,
				},
			)
			return
		} else if err.Error() == "2" {
			c.JSON(
				http.StatusForbidden,
				gin.H{
					"message": "Wrong Password",
					"status":  http.StatusForbidden,
				},
			)
			return
		}
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"message": helpers.MsgErr,
				"status":  http.StatusInternalServerError,
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"status":  http.StatusOK,
			"message": "Login Success",
			"result":  result,
		},
	)
}
