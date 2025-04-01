package utils

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/mlonV/dingtalk/types/supervisor"
)

var (
	jwtKey      = []byte("supervisor")
	tokenExpire = time.Hour * 6
)

// GenerateToken 生成 JWT token
func GenerateToken(username string) (string, error) {
	expirationTime := time.Now().Add(tokenExpire) // 设置 token 过期时间

	claims := &supervisor.Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken 验证 JWT token
func ValidateToken(tokenString string) (*supervisor.Claims, error) {
	claims := &supervisor.Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errors.New("token is malformed")
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				return nil, errors.New("token is expired")
			} else {
				return nil, errors.New("could not handle this token")
			}
		}
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func AuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		// 从请求头中获取 Authorization 字段
		err, tokenStr := GetToken(c)
		if err != nil {
			c.JSON(http.StatusOK, supervisor.Response{
				Code:    10011,
				Message: err.Error(),
				Data:    nil,
			})
			c.Abort()
			return
		}

		// 解析 token
		token, err := jwt.ParseWithClaims(tokenStr, &supervisor.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil {
			c.JSON(http.StatusOK, supervisor.Response{
				Code:    10012,
				Message: "token : " + err.Error(),
				Data:    nil,
			})
			c.Abort()
			return
		}

		// 验证 token 是否有效
		if claims, ok := token.Claims.(*supervisor.Claims); ok && token.Valid {
			c.Set("username", claims.Username)
			c.Next()
		} else {

			c.JSON(http.StatusOK, supervisor.Response{
				Code:    10013,
				Message: "Invalid token claims",
				Data:    nil,
			})
			c.Abort()
		}
	}
}

func GetToken(c *gin.Context) (error, string) {
	// 从请求头中获取 Authorization 字段
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return errors.New("Authorization header is missing"), ""
	}
	// 提取 token
	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
	if tokenString == "" {
		return errors.New("Token is missing"), ""
	}
	return nil, tokenString

}
