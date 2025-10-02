package handler

import (
	"context"
	"net/http"

	"github.com/CogitoNTNU/cogi-go/internal/service"
	"github.com/gin-gonic/gin"
)

type User struct {
	service *service.User
	ctx     *context.Context
}

func NewUserHandler(userService *service.User, ctx *context.Context) *User {
	return &User{service: userService, ctx: ctx}
}

func (u *User) GetAllUsers(gCtx *gin.Context) {
	allUsers, err := u.service.GetAllUsers(u.ctx)
	if err != nil {
		gCtx.AbortWithStatusJSON(err.Code, err.Message)
		return
	}

	gCtx.JSON(http.StatusOK, allUsers)
}

func (u *User) GetUserByID(gCtx *gin.Context) {
	userId := gCtx.Param("userId")
	user, err := u.service.GetUserByID(u.ctx, userId)

	if err != nil {
		gCtx.AbortWithStatusJSON(err.Code, err.Message)
		return
	}

	gCtx.JSON(http.StatusOK, user)
}
