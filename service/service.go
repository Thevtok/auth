package service

import (
	"github.com/Thevtok/auth/model"
	"github.com/Thevtok/auth/repo"
)

type LoginService interface {
	LoginSkuy(username string, password string) (*model.User, error)
}

type loginService struct {
	loginRepo repo.LoginRepo
}

func (u *loginService) LoginSkuy(username string, password string) (*model.User, error) {
	student, err := u.loginRepo.GetByUsernameAndPassword(username, password)
	if err != nil {
		return nil, err
	}
	if student == nil {
		return nil, nil
	}
	return student, nil
}

func NewLoginService(loginRepo repo.LoginRepo) LoginService {
	return &loginService{
		loginRepo: loginRepo,
	}
}
