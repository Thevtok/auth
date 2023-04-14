package main

import (
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtKey = []byte("thevtok")

type Student struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	// buat gin router
	router := gin.Default()

	// set up routes for login
	router.POST("/auth/login", login)

	// other routes
	studentRouter := router.Group("/students")
	studentRouter.Use(authMiddleware())

	studentRouter.GET("/", profile)

	// start server
	log.Fatal(router.Run(":8080"))
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Set("claims", claims)

		c.Next()
	}
}

func login(c *gin.Context) {
	var student Student

	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// authenticate student (compare studentname dan password)
	if student.Username == "fikri" && student.Password == "fikri" {
		// generate JWT token
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["studentname"] = student.Username
		claims["exp"] = time.Now().Add(time.Minute * 60).Unix()

		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": tokenString})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unregistered student"})
	}
}

func profile(c *gin.Context) {
	// ambil studentname dari JWT token
	claims := c.MustGet("claims").(jwt.MapClaims)
	studentname := claims["studentname"].(string)

	// dapatkan informasi student dari database (dalam hal ini, kita return studentname)
	c.JSON(http.StatusOK, gin.H{
		"message":     "Welcome to profile",
		"studentname": studentname,
	})
}
