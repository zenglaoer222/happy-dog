package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"happy-dog/utils"
	"happy-dog/utils/errmsg"
	"net/http"
	"strings"
	"time"
)

var JwtKey = []byte(utils.JwtKey)
var code int

type MyClaims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	UserId   uint   `json:"user_id"`
	jwt.StandardClaims
}

// 生成token
func SetToken(username string, role string, cid uint) (string, int) {
	expireTime := time.Now().Add(10 * time.Hour)
	//这里传入role来区分manager,customer以及shop
	setClaims := MyClaims{
		Username: username,
		Role:     role,
		UserId:   cid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "ginblog", // 签发人
		},
	}

	reqClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, setClaims)
	token, err := reqClaim.SignedString(JwtKey)
	if err != nil {
		return "", errmsg.ERROR
	}
	return token, errmsg.SUCCESS
}

// 解析token
func ParseToken(token string) (*MyClaims, error) {
	setToken, err := jwt.ParseWithClaims(token, &MyClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return JwtKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := setToken.Claims.(*MyClaims); ok && setToken.Valid {
		return claims, nil
	}

	return nil, jwt.ErrInvalidKey
}

// jwt中间件
func JwtToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenHeader := c.Request.Header.Get("Authorization")

		if tokenHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "Token is missing",
			})
			c.Abort()
			return
		}

		// 检查token格式
		tokenString := strings.SplitN(tokenHeader, "Bearer ", 2)
		if len(tokenString) != 2 {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": http.StatusBadRequest,
				"msg":  "Token format is incorrect",
			})
			c.Abort()
			return
		}

		claims, err := ParseToken(tokenString[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "Token is invalid",
			})
			c.Abort()
			return
		}

		// 检查用户角色
		role := claims.Role
		requestPath := c.Request.URL.Path

		// 定义不同角色的可访问路由
		publicRoutes := map[string]bool{
			"/api/v1/customer/add":   true,
			"/api/v1/customer/login": true,
			"/api/v1/manager/login":  true,
		}
		protectedRoutes := map[string]bool{
			"/api/v1/customer/:id":    true,
			"/api/v1/shop/:id":        true,
			"/api/v1/order/add":       true,
			"/api/v1/product/add":     true,
			"/api/v1/wallet/add":      true,
			"/api/v1/wallet/balance":  true,
			"/api/v1/order":           true,
			"/api/v1/customer/update": true,
			"/api/v1/shop":            true,
			"/api/v1/friends/add":     true,
			"/api/v1/friends":         true,
			"/api/v1/friends/wait":    true,
			"/api/v1/friends/accept":  true,
			"/api/v1/friends/search":  true,
		}
		//privateRoutes := map[string]bool{
		//	"/api/v1/customer": true,
		//	"/api/v1/shop":     true,
		//}

		// 检查角色和请求路径
		allowed := false
		if role == "manager" {
			// Role 2 可以访问所有路由
			allowed = true
		} else if role == "customer" {
			// Role 1 只能访问 protected 和 public 路由
			allowed = protectedRoutes[requestPath] || publicRoutes[requestPath]
		} else if role == "shop" {
			// Role 0 只能访问 public 路由
			allowed = publicRoutes[requestPath]
		}

		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{
				"code": http.StatusForbidden,
				"msg":  "路径权限不够，访问失败",
			})
			c.Abort()
			return
		}

		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Set("user_id", claims.UserId)
		c.Next()
	}
}
