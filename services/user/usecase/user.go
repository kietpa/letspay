package usecase

import (
	"context"
	"fmt"
	"letspay/pkg/auth"
	"letspay/pkg/logger"
	"letspay/pkg/util"
	"letspay/services/user/common/constants"
	"letspay/services/user/model"
	"letspay/services/user/repository/database"
	"net/http"
	"time"
)

func NewUserUsecase(
	userRepo database.UserRepo,
) UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
	}
}

func (u userUsecase) RegisterUser(
	ctx context.Context, registerUserRequest model.RegisterUserRequest,
) (model.UserDetail, model.Error) {
	_, err := u.userRepo.GetUserByEmail(ctx, registerUserRequest.Email)
	if err == nil {
		logger.Error(ctx, fmt.Sprintf("[Register User] check user email conflict error=%s email=%s",
			err,
			registerUserRequest.Email,
		))
		return model.UserDetail{}, model.Error{
			Code:    http.StatusConflict,
			Message: constants.INVALID_EMAIL_MESSAGE,
		}
	}

	hashedPass, err := util.HashPassword(registerUserRequest.Password)
	if err != nil {
		logger.Error(ctx, fmt.Sprintf("[Register User] password hash error=%s email=%s",
			err,
			registerUserRequest.Email,
		))
		return model.UserDetail{}, model.Error{
			Code:    http.StatusInternalServerError,
			Message: constants.INTERNAL_ERROR_MESSAGE,
		}
	}

	err = u.userRepo.RegisterUser(ctx, model.RegisterUserInput{
		Name:           registerUserRequest.Name,
		Email:          registerUserRequest.Email,
		HashedPassword: hashedPass,
		CreatedAt:      time.Now(),
	})
	if err != nil {
		logger.Error(ctx, fmt.Sprintf("[Register User] repo error=%s email=%s",
			err,
			registerUserRequest.Email,
		))
		return model.UserDetail{}, model.Error{
			Code:    http.StatusInternalServerError,
			Message: constants.INTERNAL_ERROR_MESSAGE,
		}
	}

	return model.UserDetail{
		Name:  registerUserRequest.Name,
		Email: registerUserRequest.Email,
	}, model.Error{}
}

func (u userUsecase) LoginUser(
	ctx context.Context, loginUserRequest model.LoginUserRequest,
) (model.LoginUserResponse, model.Error) {

	user, err := u.userRepo.GetUserByEmail(ctx, loginUserRequest.Email)
	if err != nil {
		logger.Error(ctx, fmt.Sprintf("[Login User] check user DB error=%s email=%s",
			err,
			loginUserRequest.Email,
		))
		return model.LoginUserResponse{}, model.Error{
			Code:    http.StatusNotFound,
			Message: constants.INVALID_EMAIL_MESSAGE,
		}
	}

	if !util.CheckPassword(loginUserRequest.Password, user.HashedPassword) {
		logger.Error(ctx, fmt.Sprintf("[Login User] invalid password error=%s email=%s",
			err,
			loginUserRequest.Email,
		))
		return model.LoginUserResponse{}, model.Error{
			Code:    http.StatusBadRequest,
			Message: constants.INVALID_EMAIL_MESSAGE,
		}
	}

	token, err := auth.GenerateToken(user.UserId)
	if err != nil {
		logger.Error(ctx, fmt.Sprintf("[Login User] generate token error=%s email=%s userid=%d",
			err,
			loginUserRequest.Email,
			user.UserId,
		))
		return model.LoginUserResponse{}, model.Error{
			Code:    http.StatusInternalServerError,
			Message: constants.INTERNAL_ERROR_MESSAGE,
		}
	}

	return model.LoginUserResponse{
		Token: token,
	}, model.Error{}
}

func (u userUsecase) GetUser(
	ctx context.Context,
	userId int,
) (model.GetUserDetail, model.Error) {
	user, err := u.userRepo.GetUserById(ctx, userId)
	if err != nil {
		logger.Error(ctx, fmt.Sprintf("[Get User] Get user DB error=%s id=%d",
			err,
			userId,
		))
		return model.GetUserDetail{}, model.Error{
			Code:    http.StatusNotFound,
			Message: "User id not found",
		}
	}

	return model.GetUserDetail{
		UserId:  userId,
		Name:    user.Name,
		Email:   user.Email,
		Webhook: user.Webhook,
	}, model.Error{}
}

func (u userUsecase) AddWebhook(
	ctx context.Context,
	input model.AddWebhookRequest,
) model.Error {
	err := u.userRepo.UpdateUserWebhook(ctx, input)
	if err != nil {
		logger.Error(ctx, fmt.Sprintf(
			"[Add Webhook] update disbursement DB error=%s refid=%d",
			err,
			input.UserId,
		))
		return model.Error{
			Code:    http.StatusInternalServerError,
			Message: constants.INTERNAL_ERROR_MESSAGE,
		}
	}

	return model.Error{}
}
