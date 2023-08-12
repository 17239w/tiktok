package userlogin

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
