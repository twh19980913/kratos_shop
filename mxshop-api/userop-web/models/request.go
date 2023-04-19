package models

import (
	"github.com/dgrijalva/jwt-go"
)

type CustomClaims struct {
	ID          uint // 用户ID
	NickName    string // 用户名称
	AuthorityId uint // 对应用户权限
	jwt.StandardClaims
}
