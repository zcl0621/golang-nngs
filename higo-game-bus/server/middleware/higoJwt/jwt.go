package higoJwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NeedJwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":     -1,
				"error_code": "请求未携带token，无权限访问",
			})
			c.Abort()
			return
		}
		j := NewJWT()
		// parseToken 解析token包含的信息
		claims, err := j.ParseToken(token)
		if err != nil {
			if err == TokenExpired {
				c.JSON(http.StatusUnauthorized, gin.H{
					"status":     -1,
					"error_code": "授权已过期",
				})
				c.Abort()
				return
			}
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":     -1,
				"error_code": err.Error(),
			})
			c.Abort()
			return
		}
		c.Set("claims", claims)
	}
}

// JWT 签名结构
type JWT struct {
	SigningKey []byte
}

var (
	TokenExpired       = errors.New("token is expired")
	TokenNotValidYet   = errors.New("token not active yet")
	TokenMalformed     = errors.New("that's not even a token")
	TokenInvalid       = errors.New("couldn't handle this token:")
	SignKey           = "eStarGo"
)

type CustomClaims struct {
	ID            string `json:"user_id"`
	StudentName   string `json:"student_name"`
	StudentMobile string `json:"student_mobile"`

	jwt.StandardClaims
}

func NewJWT() *JWT {
	return &JWT{
		[]byte(GetSignKey()),
	}
}

func GetSignKey() string {
	return SignKey
}

func SetSignKey(key string) string {
	SignKey = key
	return SignKey
}

func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, TokenInvalid
}

func GetTokenInfo(c *gin.Context) (*CustomClaims, error) {
	if value, hasValue := c.Get("claims"); !hasValue {
		return nil, errors.New("没有token")
	} else {
		if data, ok := value.(*CustomClaims); ok {
			return data, nil
		} else {
			return nil, errors.New("token 解析错误")
		}
	}
}
