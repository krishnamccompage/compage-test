package services

import (
	"github.com/krishnamccompage/compage-test/test-api/pkg/rest/server/daos"
	"github.com/krishnamccompage/compage-test/test-api/pkg/rest/server/models"
)

type UserService struct {
	userDao *daos.UserDao
}

func NewUserService() (*UserService, error) {
	userDao, err := daos.NewUserDao()
	if err != nil {
		return nil, err
	}
	return &UserService{
		userDao: userDao,
	}, nil
}

func (userService *UserService) CreateUser(user *models.User) (*models.User, error) {
	return userService.userDao.CreateUser(user)
}

func (userService *UserService) UpdateUser(id string, user *models.User) (*models.User, error) {
	return userService.userDao.UpdateUser(id, user)
}

func (userService *UserService) DeleteUser(id string) error {
	return userService.userDao.DeleteUser(id)
}

func (userService *UserService) ListUsers() ([]*models.User, error) {
	return userService.userDao.ListUsers()
}

func (userService *UserService) GetUser(id string) (*models.User, error) {
	return userService.userDao.GetUser(id)
}
