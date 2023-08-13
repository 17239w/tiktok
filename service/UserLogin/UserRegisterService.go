package userlogin

import (
	"errors"
	"fmt"
	"log"
	"tiktok/middleware"
	"tiktok/models"
)

// UserRegisterService：用户注册服务
type UserRegisterService struct {
	username string
	password string

	userloginresponse *UserLoginResponse // 保存 用户登陆成功后的响应数据，包括token和id
	userid            int64
	token             string
}

// UserRegister：注册用户，得到token和id
func UserRegister(username, password string) (*UserLoginResponse, error) {
	//打印日志
	log.Printf("Service层(注册用户，得到token和id)：username:%s, password:%s", username, password)
	return NewUserRegisterService(username, password).Do()
}

// NewUserRegisterService：创建一个UserRegisterService结构体
func NewUserRegisterService(username, password string) *UserRegisterService {
	return &UserRegisterService{username: username, password: password}
}

// Do：执行注册服务
func (service *UserRegisterService) Do() (*UserLoginResponse, error) {
	//1.对参数进行合法性验证
	if err := service.checkNum(); err != nil {
		return nil, err
	}
	//2.更新数据到数据库
	if err := service.updateData(); err != nil {
		return nil, err
	}
	//3.打包response
	if err := service.packResponse(); err != nil {
		return nil, err
	}
	return service.userloginresponse, nil
}

// 1.对参数进行合法性验证
func (service *UserRegisterService) checkNum() error {
	//打印日志
	log.Printf("Service层(对参数进行合法性验证)：username:%s, password:%s", service.username, service.password)
	if service.username == "" {
		fmt.Println("username is null")
		return errors.New("username is null")
	}
	if len(service.username) > MaxUsernameLength {
		fmt.Println("the length of username is too long")
		return errors.New("the length of username is too long")
	}

	if service.password == "" {
		fmt.Println("password is null")
		return errors.New("password is null")
	}
	return nil
}

// 2.更新数据到数据库
func (service *UserRegisterService) updateData() error {

	//打印日志
	log.Printf("models层:(将login信息保存到数据库)：username:%s, password:%s", service.username, service.password)

	//调用models层，创建一个userLogin结构体和一个userinfo结构体
	userLogin := models.UserLogin{Username: service.username, Password: service.password}
	userinfo := models.UserInfo{User: &userLogin, Name: service.username}

	//调用models层，判断用户名是否已经存在
	userLoginDAO := models.NewUserLoginDao()
	if userLoginDAO.IsUserExist(service.username) {
		return errors.New("用户名已存在")
	}

	//调用models层，更新操作 (由于userLogin的父表是userInfo，故更新userInfo即可)
	userInfoDAO := models.NewUserInfoDAO()
	err := userInfoDAO.AddUserInfo(&userinfo)
	if err != nil {
		return err
	}

	//调用middleware层，颁发token
	token, err := middleware.ReleaseToken(userLogin)
	if err != nil {
		return err
	}
	service.token = token
	service.userid = userinfo.Id
	return nil
}

// 3.打包UserLoginResponse
func (service *UserRegisterService) packResponse() error {
	//将userid和token打包到response中
	service.userloginresponse = &UserLoginResponse{
		UserId: service.userid,
		Token:  service.token,
	}
	return nil
}
