package usecase

import (
	"context"
	"fmt"
	"letspay/common/constants"
	"letspay/model"
	"letspay/repository/database"
	"letspay/tool/logger"
	"letspay/tool/util"
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
	//check if user exists
	_, err := u.userRepo.GetUserByEmail(ctx, registerUserRequest.Email)
	if err == nil {
		logger.Error(ctx, fmt.Sprintf("[Register User - Usecase] check user email conflict: %s email=%s",
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
		logger.Error(ctx, fmt.Sprintf("[Register User - Usecase] password hash error: %s email=%s",
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
		logger.Error(ctx, fmt.Sprintf("[Register User - Usecase] repo error: %s email=%s",
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
		logger.Error(ctx, fmt.Sprintf("[Login User - Usecase] check user DB error: %s email=%s",
			err,
			loginUserRequest.Email,
		))
		return model.LoginUserResponse{}, model.Error{
			Code:    http.StatusNotFound,
			Message: constants.INVALID_EMAIL_MESSAGE,
		}
	}

	if !util.CheckPassword(loginUserRequest.Password, user.HashedPassword) {
		logger.Error(ctx, fmt.Sprintf("[Login User - Usecase] invalid password error: %s email=%s",
			err,
			loginUserRequest.Email,
		))
		return model.LoginUserResponse{}, model.Error{
			Code:    http.StatusBadRequest,
			Message: constants.INVALID_EMAIL_MESSAGE,
		}
	}

	token, err := util.GenerateToken(&user)
	if err != nil {
		logger.Error(ctx, fmt.Sprintf("[Login User - Usecase] generate token error: %s email=%s userid=%d",
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
