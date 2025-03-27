package super

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/mlonV/dingtalk/config"
	"github.com/mlonV/dingtalk/types/supervisor"
	"github.com/mlonV/dingtalk/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

var (
	Mysql = &config.Conf.Mysql
	dsn   = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		Mysql.Username,
		Mysql.Password,
		Mysql.Hostname,
		Mysql.Port,
		Mysql.Database)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
)

func Login(c *gin.Context) {

	var user supervisor.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dbUser, err := getUserByName(user.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if dbUser.Username == "" {
		c.JSON(http.StatusOK, supervisor.Response{
			Code:    10019,
			Data:    dbUser,
			Message: "user not found",
		})
		return
	}

	if user.Password != dbUser.Password {
		c.JSON(http.StatusOK, supervisor.Response{
			Code:    10019,
			Data:    dbUser,
			Message: "密码错误，请重试",
		})
		return
	}

	token, err := utils.GenerateToken(dbUser.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, supervisor.Response{
			Code:    10019,
			Data:    nil,
			Message: "token生成失败",
		})
		return
	}
	c.JSON(http.StatusOK, supervisor.Response{
		Code:    0,
		Data:    supervisor.ClientToken{Token: token},
		Message: "login success",
	})
}

func UserInfo(c *gin.Context) {
	var userinfo supervisor.UserInfo
	// 从请求头中获取 Authorization 字段
	err, tokenStr := utils.GetToken(c)
	if err != nil {
		c.JSON(http.StatusOK, supervisor.Response{
			Code:    10019,
			Message: err.Error(),
			Data:    nil,
		})
		c.Abort()
		return
	}
	claims, err := utils.ValidateToken(tokenStr)
	if err != nil {
		c.JSON(http.StatusOK, supervisor.Response{
			Code:    10019,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}

	user, err := getUserByName(claims.Username)
	userinfo.Avatar = user.Avatar
	userinfo.Introduction = user.Introduction
	userinfo.Name = user.Username
	userinfo.Roles = []string{user.Username}
	// userinfo.Token = token
	c.JSON(http.StatusOK, supervisor.Response{
		Code:    0,
		Data:    userinfo,
		Message: "success",
	})
}

func LogOut(c *gin.Context) {

	// 从请求头中获取 Authorization 字段
	err, tokenStr := utils.GetToken(c)
	if err != nil {
		c.JSON(http.StatusOK, supervisor.Response{
			Code:    10019,
			Message: err.Error(),
			Data:    nil,
		})
		c.Abort()
		return
	}
	claims, err := utils.ValidateToken(tokenStr)
	if err != nil {
		c.JSON(http.StatusOK, supervisor.Response{
			Code:    10019,
			Data:    nil,
			Message: err.Error(),
		})
		c.Abort()
		return
	}
	claims.ExpiresAt = time.Now().Unix()

	c.JSON(http.StatusOK, supervisor.Response{
		Code:    0,
		Data:    nil,
		Message: "Logout success",
	})
}

func getUserByName(username string) (*supervisor.User, error) {
	var dbUser supervisor.User
	if err := db.Where("username = ? ", username).First(&dbUser).Error; err != nil {
		fmt.Println(err, dbUser)
	}
	return &dbUser, err
}

func addUser(username, password string) error {
	var user = supervisor.User{Username: username, Password: password}
	result := db.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
