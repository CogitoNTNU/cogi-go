package service

import (
	"context"
	"net/http"

	"github.com/CogitoNTNU/cogi-go/internal/model"
	userRepository "github.com/CogitoNTNU/cogi-go/internal/repository/user"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type User struct {
	logger     *logrus.Entry
	repository *userRepository.Repo
}

func NewUserService(userRepository *userRepository.Repo, logger *logrus.Entry) *User {
	return &User{repository: userRepository, logger: logger}
}

func (u *User) GetAllUsers(ctx *context.Context) ([]model.User, *model.ErrorResponse) {
	allUsers, err := u.repository.GetAllUsers(ctx)
	if err != nil {
		u.logger.Errorf("An error has occured when retrieving all users. Error code: %d", err.Code)
		return nil, err
	}

	return allUsers, nil
}

func (u *User) GetUserByID(ctx *context.Context, userId string) (*model.User, *model.ErrorResponse) {
	id, err := uuid.Parse(userId)
	if err != nil {
		return nil, &model.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	user, fetchErr := u.repository.GetUserByID(ctx, id)
	if fetchErr != nil {
		return nil, fetchErr
	}

	return user, nil
}
