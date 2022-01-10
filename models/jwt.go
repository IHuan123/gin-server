package models

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var jwtKey = []byte("admin.com")

type Claims struct{
	Uid int
	UserName string
	Password string
	LoginTime int64
	jwt.StandardClaims
}

//生产token
func GenerateToken(Uid int,userName string,phone string)(string,error){
	expireTime := time.Now().Add(7 * 24 * 3600 * time.Second).Unix()
	claims := &Claims{
		Uid: Uid,
		UserName: userName,
		Password: phone,
		LoginTime: time.Now().Unix(),
		StandardClaims:jwt.StandardClaims{
			ExpiresAt: expireTime, //过期时间
			IssuedAt:  time.Now().Unix(),
			Issuer:    "192.168.0.100",  // 签名颁发者
			Subject:   "user token", //签名主题
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	tokenString,err := token.SignedString(jwtKey)
	if err != nil {
		fmt.Println(err)
		return "",err
	}
	return tokenString,nil
}
//解析token
func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (i interface{}, err error) {
		return jwtKey, nil
	})
	return token, claims, err
}

