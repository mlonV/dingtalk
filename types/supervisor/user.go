package supervisor

import "github.com/dgrijalva/jwt-go"

type User struct {
	ID           int    `gorm:"primaryKey;table:user"`
	UserID       int    `gorm:"column:user_id"`
	Username     string `gorm:"column:username"  json:"username"`
	Password     string `gorm:"column:password" json:"password"`
	Avatar       string `gorm:"column:avatar" json:"avatar"`
	Introduction string `gorm:"column:introduction" json:"introduction"`
}

type ClientToken struct {
	Token string `json:"token"`
}

type UserInfo struct {
	Roles        []string `json:"roles"`
	Introduction string   `json:"introduction"`
	Avatar       string   `json:"avatar"`
	Name         string   `json:"name"`
	// Token        string   `json:"token"`
}

// 明确指定表名为 user
func (User) TableName() string {
	return "user"
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Type    string      `json:"type"`
	Data    interface{} `json:"data"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
