package util

import (
	"Cloud/pkg/consts"
	"time"

	"Cloud/pkg/setting"
	"github.com/dgrijalva/jwt-go"
)

var jwtSecret = []byte(setting.JwtSecret)

type Claims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Authority int   `json:"authority"`
	jwt.StandardClaims
}

func GenerateToken(username,password string , authority int) (string,error){
	nowTime := time.Now()
	expireTime := nowTime.Add(consts.EXPIRE_DURATION)

	claims := Claims{
		username,
		password,
		authority,
		jwt.StandardClaims{
			ExpiresAt:expireTime.Unix(),
			Issuer:"apiTest",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	token,err := tokenClaims.SignedString(jwtSecret)

	return token,err
}

func ParseToken(token string)(*Claims,error){
	tokenClaims, err := jwt.ParseWithClaims(token,&Claims{}, func(token *jwt.Token) (i interface{}, e error) {
		return jwtSecret,nil
	})

	if tokenClaims != nil{
		if claims,ok := tokenClaims.Claims.(*Claims);ok && tokenClaims.Valid{
			return claims,nil
		}
	}

	return nil,err
}

//func GetUserAuth(token string)int{
//	claims,err := ParseToken(token)
//
//	if err != nil{
//		log.Fatal(err)
//		return 0
//	}
//
//	return claims.Authority
//}

