package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/jwt"
	"bluebell/pkg/snowflake"
)

// SignUp 处理注册业务相关逻辑
func SignUp(p *models.ParamSignUp) (tokenerr error) {
	// 1. 判断用户是否存在
	if err := mysql.CheckUserExist(p.Username); err != nil {
		return err
	}
	// 2. 生成UID
	userID := snowflake.GenID()
	// 构造一个User实例
	user := &models.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}
	// 3. 保存进数据库
	if err := mysql.InsertUser(user); err != nil {
		return err
	}
	return nil
}

// Login 处理登录业务相关逻辑
func Login(p *models.ParamLogin) (token string, err error) {
	user := &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	// 传递的是指针，就能拿到user.UserID
	if err = mysql.Login(user); err != nil {
		return "", err
	}
	return jwt.GenToken(user.UserID)
}
