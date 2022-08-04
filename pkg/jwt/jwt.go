package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

//自定义声明类型 并内嵌jwt.RegisteredClaims
//注意大写可访问性
type MyClaims struct {
	Username string `json:"username"`
	//Password string `json:"password"`
	jwt.RegisteredClaims
}

// CustomSecret 用于加盐的字符串
var CustomSecret = []byte("coco")

// TokeExpireDuration 过期时间
const TokeExpireDuration = time.Hour * 1000

// GenRegisteredClaims 生成jwt
func GenRegisteredClaims(u string) (string, error) {

	claims := MyClaims{
		Username: u,
		//Password: p,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "coco",                                                 //签发人
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokeExpireDuration)), //过期时间
		},
	}
	//生成Token对象
	//NewWithClaims使用指定的签名方法和声明创建新令牌。
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//生成签名字符串
	//SignedString创建并返回一个完整的带签名的JWT。
	//使用令牌中指定的SigningMethod对令牌进行签名。
	return token.SignedString(CustomSecret)
}

//解析jwt
func ParseToken(tokenString string) (*MyClaims, error) {
	//解析Token
	//自定义Claim结构体则需要使用 ParseWithClaims 方法
	//Token表示JWT令牌。
	//根据您是在创建令牌还是解析/验证令牌，将使用不同的字段。
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return CustomSecret, nil
	})
	if err != nil {
		return nil, err
	}

	// 对token对象中的Claim进行类型断言
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid { // 校验token
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
