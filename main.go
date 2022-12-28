package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"example.com/models"
	"example.com/controllers"
	"github.com/appleboy/gin-jwt/v2"
	"time"
	"log"
)

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

var identityKey = "id"

// User demo
type User struct {
	UserName  string
	FirstName string
	LastName  string
}

func main() {

	r := gin.Default();

	// the jwt middleware
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "example zone",
		Key:         []byte("harusdiacakbiarsusah"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
		  if v, ok := data.(*User); ok {
			return jwt.MapClaims{
			  identityKey: v.UserName,
			}
		  }
		  return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
		  claims := jwt.ExtractClaims(c)
		  return &User{
			UserName: claims[identityKey].(string),
		  }
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
		  var loginVals login
		  if err := c.ShouldBind(&loginVals); err != nil {
			return "", jwt.ErrMissingLoginValues
		  }
		  userID := loginVals.Username
		  password := loginVals.Password
	
		  if (userID == "admin" && password == "admin") || (userID == "test" && password == "test") {
			return &User{
			  UserName:  userID,
			  LastName:  "Bo-Yi",
			  FirstName: "Wu",
			}, nil
		  }
	
		  return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
		  if v, ok := data.(*User); ok && v.UserName == "admin" {
			return true
		  }
	
		  return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
		  c.JSON(code, gin.H{
			"code":    code,
			"message": message,
		  })
		},
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
	
		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	  })
	
	  if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	  }
	  errInit := authMiddleware.MiddlewareInit()
	
	  if errInit != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	  }
	
	  r.POST("/login", authMiddleware.LoginHandler)
	
	  r.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	  })

	//MODEL
	db := models.SetupModels()

	r.Use(func(c *gin.Context){
		c.Set("db", db)
		c.Next()
	})
	
	  auth := r.Group("/auth")
	  // Refresh time can be longer than token timeout
	  auth.GET("/refresh_token", authMiddleware.RefreshHandler)
	  auth.Use(authMiddleware.MiddlewareFunc())
	  {
		auth.GET("/hello", func(c *gin.Context){
			c.JSON(http.StatusOK, gin.H{"data":"Berhasil Login"})
		})

		auth.GET("/nasabah", controllers.NasabahTampil)
		auth.POST("/nasabah", controllers.NasabahTambah)
		auth.PUT("/nasabah/:nim", controllers.NasabahUbah)
		auth.DELETE("/nasabah/:nim", controllers.NasabahHapus)
	  }


	r.GET("/", func(c *gin.Context){
		c.JSON(http.StatusOK, gin.H{"data":"pong"})
	})

	r.Use(func(c *gin.Context){
		c.Set("db", db)
		c.Next()
	})

	r.Run()
}