package main

import (
	"time"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
)

type User struct {
	Username string
	Password string
}

type Todo struct {
	Deskripsi string
}

var identityKey = "username"

func main() {
	router := gin.Default()

	authMiddleware, _ := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("dscunikomgolang"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					identityKey: v.Username,
					"apaya":     "unikom",
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &Todo{
				Deskripsi: claims[identityKey].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals User
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			userID := loginVals.Username
			password := loginVals.Password

			/**
			Ngambil data dari database
			*/

			if userID == "admin" && password == "admin" {
				return &User{
					Username: "admin",
					Password: "admin",
				}, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if v, ok := data.(*User); ok && v.Username == "admin" {
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"Kode":  code,
				"Pesan": message,
			})
		},
		TokenLookup:   "header: Authorization, query: authorization, cookie: jwt",
		TokenHeadName: "James",
		TimeFunc:      time.Now,
	})

	router.POST("/login", authMiddleware.LoginHandler)
	coba := router.Group("/user")
	coba.Use(authMiddleware.MiddlewareFunc())
	{
		coba.GET("/", helloHandler)
	}
	router.Run()
}

func helloHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	user, _ := c.Get(identityKey)
	c.JSON(200, gin.H{
		"userID":   claims["username"],
		"apaya":    claims["apaya"],
		"userName": user.(*Todo).Deskripsi,
		"text":     "Hello World.",
	})
}
