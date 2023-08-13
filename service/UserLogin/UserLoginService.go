package userlogin

import (
	"errors"
	"tiktok/middleware"
	"tiktok/models"
)

const (
	MaxUsernameLength = 10 // 用户名最大长度
	MaxPasswordLength = 20 // 密码最大长度
	MinPasswordLength = 8  // 密码最小长度
)

// UserLoginResponse：用户登录响应
type UserLoginResponse struct {
	UserId int64  `json:"user_id"`
	Token  string `json:"token"`
}

// UserLoginService：用户登录服务
type UserLoginService struct {
	username string
	password string

	data   *UserLoginResponse // 用户登录响应
	userid int64
	token  string
}

// QueryUserLogin:查询用户是否存在，并返回token和id
func QueryUserLogin(username, password string) (*UserLoginResponse, error) {
	return NewUserLoginService(username, password).Do()
}

// NewUserLoginService:创建用户登录服务
func NewUserLoginService(username, password string) *UserLoginService {
	return &UserLoginService{username: username, password: password}
}

// Do：执行查询流程
func (service *UserLoginService) Do() (*UserLoginResponse, error) {
	//1.对参数进行合法性验证
	if err := service.checkNum(); err != nil {
		return nil, err
	}
	//2.准备好数据
	if err := service.prepareData(); err != nil {
		return nil, err
	}
	//3.打包最终数据
	if err := service.packData(); err != nil {
		return nil, err
	}
	return service.data, nil
}

// 1.对参数进行合法性验证
func (service *UserLoginService) checkNum() error {
	if service.username == "" {
		return errors.New("用户名为空")
	}
	if len(service.username) > MaxUsernameLength {
		return errors.New("用户名长度超出限制")
	}
	if service.password == "" {
		return errors.New("密码为空")
	}
	return nil
}

// 2.准备好数据
func (service *UserLoginService) prepareData() error {
	userLoginDAO := models.NewUserLoginDao()
	var login models.UserLogin
	//调用models层，查询用户是否存在
	err := userLoginDAO.QueryUserLogin(service.username, service.password, &login)
	if err != nil {
		return err
	}
	service.userid = login.UserInfoId
	//颁发token
	token, err := middleware.ReleaseToken(login)
	if err != nil {
		return err
	}
	service.token = token
	return nil
}

// 3.打包最终数据
func (service *UserLoginService) packData() error {
	service.data = &UserLoginResponse{
		UserId: service.userid,
		Token:  service.token,
	}
	return nil
}
