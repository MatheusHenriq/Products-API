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

func (uu *UserUsecase) CreateUser(user model.User) error {
	err := uu.repository.CreateUser(user)
	if err != nil {
		return err
	}
	return nil
}

func (uu *UserUsecase) DeleteUser(user model.User) (*model.User, error) {
	userData, err := uu.repository.DeleteUser(user)
	if err != nil {
		return userData, err
	}
	return userData, nil
}

func (uu *UserUsecase) LogIn(user model.User) (*model.User, string, error) {
	userData, token, err := uu.repository.LogIn(user)
	if err != nil {
		return userData, token, err
	}
	return userData, token, nil
}
