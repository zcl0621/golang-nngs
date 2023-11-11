package adminJwt

import (
	"errors"
	"fmt"
	"higo-game-bus/redisUtils"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
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
		var tokenStr = fmt.Sprintf("adminToken_%d", claims.ID)
		cacheToken, _ := redisUtils.Get(tokenStr)
		if fmt.Sprintf("%s", cacheToken) != token {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": -2,
				"errlr{": "已在其他设备上登录，请重新登录",
			})
			c.Abort()
			return
		}
		// 继续交由下一个路由处理,并将解析出的信息传递下去
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
	SignKey           = "eStarGoAdmin"
)

type CustomClaims struct {
	ID              uint   `json:"user_id"`
	SchoolID        int    `json:"school_id"`
	BranchID        int    `json:"branch_id"`
	TeacherID       int    `json:"teacher_id"`
	TeacherName     string `json:"teacher_name"`
	HeadTeacherID   int    `json:"head_teacher_id"`
	HeadTeacherName string `json:"head_teacher_name"`
	IsSuperUser     bool   `json:"is_super_user"`
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

func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
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

func (j *JWT) RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return j.CreateToken(*claims)
	}
	return "", TokenInvalid
}

func FromTokenGetAdmin(token string) (adminID uint, schoolID int, teacherID int, teacherName string, headTeacherID int, headTeacherName string, isSuperUser bool, error error) {
	// 从token中获取用户id
	jwt := NewJWT()
	if claims, err := jwt.ParseToken(token); err != nil {
		return 0, 0, 0, "", 0, "", false, err
	} else {
		return claims.ID, claims.SchoolID, claims.TeacherID, claims.TeacherName, claims.HeadTeacherID, claims.HeadTeacherName, claims.IsSuperUser, nil
	}
}
