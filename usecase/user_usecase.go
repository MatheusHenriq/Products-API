package usecase

import (
	"go-api/model"
	"go-api/repository"
)

type UserUsecase struct {
	repository repository.UserRepository
}

func NewUserUsecase(repo repository.UserRepository) UserUsecase {
	return UserUsecase{repository: repo}
}

func (uu *UserUsecase) CreateUser(user model.User) (model.User, error) {
	userId, err := uu.repository.CreateUser(user)
	if err != nil {
		return model.User{}, err
	}
	user.ID = userId
	return user, nil
}

func (uu *UserUsecase) DeleteUser(user model.User) (*model.User, error) {
	userData, err := uu.repository.DeleteUser(user)
	if err != nil {
		return userData, err
	}
	return userData, nil
}
