package usecase

import (
	"context"
	"letspay/common/constants"
	"letspay/model"
	"letspay/repository/database"
	"letspay/tool/util"
	"log"
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
		log.Println("usecase register user check user err:", err)
		return model.UserDetail{}, model.Error{
			Code:    http.StatusBadRequest,
			Message: constants.EMAIL_INVALID_MESSAGE,
		}
	}

	hashedPass, err := util.HashPassword(registerUserRequest.Password)
	if err != nil {
		log.Println("usecase register user hash err:", err)
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
		log.Println("usecase register user register user err:", err)
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
