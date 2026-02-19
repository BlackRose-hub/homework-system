package service

import (
	"homework-system/dao"
	"homework-system/models"
	"homework-system/pkg/errcode"
	"homework-system/pkg/jwt"
)

type UserService struct {
	userDao *dao.UserDao
}

func NewUserService() *UserService {
	return &UserService{
		userDao: dao.NewUserDao(),
	}
}

type RegisterRequest struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	Nickname   string `json:"nickname" binding:"required"`
	Department string `json:"department" binding:"required"`
}
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type LoginResponse struct {
	AccessToken  string                 `json:"access_token"`
	RefreshToken string                 `json:"refresh_token"`
	User         map[string]interface{} `json:"user"`
}

func (s *UserService) Register(req *RegisterRequest) (*models.User, int, error) {
	exists, err := s.userDao.CheckUsernameExist(req.Username)
	if err != nil {
		return nil, errcode.ServerError, err
	}
	if exists {
		return nil, errcode.UserAlreadyExists, nil
	}
	user := &models.User{
		Username:   req.Username,
		Password:   req.Password,
		Nickname:   req.Nickname,
		Department: models.Department(req.Department),
		Role:       models.RoleStudent,
	}
	if err := user.HashPassword(); err != nil {
		return nil, errcode.ServerError, err
	}
	if err := s.userDao.Create(user); err != nil {
		return nil, errcode.ServerError, err
	}
	return user, errcode.Success, err
}
func (s *UserService) Login(req *LoginRequest) (*LoginResponse, int, error) {
	user, err := s.userDao.FindByUsername(req.Username)
	if err != nil {
		return nil, errcode.UserNotFound, err
	}
	if !user.CheckPassword(req.Password) {
		return nil, errcode.PassWordIncorrect, nil
	}
	accessToken, refreshToken, err := jwt.GenerateTokens(user.ID, string(user.Role))
	if err != nil {
		return nil, errcode.ServerError, err
	}
	response := &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: map[string]interface{}{
			"id":               user.ID,
			"Username":         user.Username,
			"Nickname":         user.Nickname,
			"role":             user.Role,
			"department":       user.Department,
			"department_label": models.DepartmentLabel[user.Department],
			"email":            user.Email,
		},
	}
	return response, errcode.Success, nil
}
func (s *UserService) GetProfile(userID uint) (*models.User, int, error) {
	user, err := s.userDao.FindByID(userID)
	if err != nil {
		return nil, errcode.UserNotFound, err
	}
	return user, errcode.Success, nil
}

func (s *UserService) DeleteAccount(userID uint) (int, error) {
	if err := s.userDao.Delete(userID); err != nil {
		return errcode.ServerError, err
	}
	return errcode.Success, nil
}

func (s *UserService) RefreshToken(refreshToken string) (string, string, int, error) {
	newAccessToken, newRefreshToken, err := jwt.RefreshTokens(refreshToken)
	if err != nil {
		return "", "", errcode.TokenInvalid, err
	}
	return newAccessToken, newRefreshToken, errcode.Success, nil
}
