/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    jwt
 * @Date:    2021/6/21 11:45 上午
 * @package: jwt
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package jwt

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jager/hawox/errcode"
	"github.com/jager/hawox/httpx"
	"time"
)

var opt = &Option{
	TokenHeaderKey: "X-Token",
	Secret:         "64981e500279991b18c4ca082fc18d44",
	Timeout:        time.Hour * 12,
}

type Option struct {
	TokenHeaderKey string
	Secret         string
	Timeout        time.Duration
}

func SetOption(opfs ...func(opt *Option)) {
	for _, opf := range opfs {
		opf(opt)
	}
}

type Claims struct {
	Uid string `json:"uid"`
	jwt.StandardClaims
}

func GenerateToken(uid string) (string, error) {
	return GenerateTokenWithTimeout(uid, opt.Timeout)
}

func GenerateTokenWithTimeout(uid string, timeout time.Duration) (string, error) {
	expireTime := time.Now().Add(timeout).Unix()
	claims := Claims{
		uid,
		jwt.StandardClaims{
			ExpiresAt: expireTime,
			Issuer:    "HawOs",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString([]byte(opt.Secret))

	return token, err
}

func ParseToken(token string) (*Claims, error) {
	claims := &Claims{}
	tokenClaims, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(opt.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !tokenClaims.Valid {
		return nil, claims.Valid()
	}

	return claims, nil
}

func RefreshTokenByToken(token string, timeout ...time.Duration) (newToken string, err error) {
	claims, err := ParseToken(token)
	if err != nil {
		return
	}
	if len(timeout) > 0 {
		newToken, err = GenerateTokenWithTimeout(claims.Uid, timeout[0])
	} else {
		newToken, err = GenerateToken(claims.Uid)
	}

	return
}

func RefreshToken(c *gin.Context) (string, bool) {
	uid, ok := Uid(c)
	if !ok {
		return "", false
	}
	token, err := GenerateToken(uid)
	if err != nil {
		httpx.ErrInterrupt(c, errcode.InternalErr.WithErr(err))
		return "", false
	}
	return token, true
}

func CheckToken(c *gin.Context) {
	token := c.GetHeader("X-Token")
	claims, err := ParseToken(token)
	if err != nil {
		httpx.ErrInterrupt(c, errcode.VerifyErr)
		return
	}
	//if time.Now().Unix() > claims.ExpiresAt {
	//	httpx.ErrInterrupt(c, errcode.VerifyErr.WithMsg("token过期了"))
	//	return
	//}
	c.Set("uid", claims.Uid)
	c.Next()
}

func Uid(c *gin.Context) (string, bool) {
	uid, ok := c.Get("uid")
	if !ok {
		httpx.ErrInterrupt(c, errcode.VerifyErr)
		return "", false
	}
	return uid.(string), true
}
