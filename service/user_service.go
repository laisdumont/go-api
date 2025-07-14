package service

import (
	"go-api/model"
	"go-api/repository"
)

type UserService struct {
	Repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{Repo: repo}
}

func (s *UserService) Create(user *model.User) error {
	return s.Repo.Create(user)
}

func (s *UserService) GetAll() ([]model.User, error) {
	return s.Repo.GetAll()
}

func (s *UserService) Update(user *model.User) error {
	return s.Repo.Update(user)
}

func (s *UserService) Delete(id int) error {
	return s.Repo.Delete(id)
}

func (s *UserService) FindByName(name string) (*model.User, error) {
	return s.Repo.FindByName(name)
}
