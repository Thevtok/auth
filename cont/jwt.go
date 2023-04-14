package cont

import (
	"net/http"
	"time"

	"github.com/Thevtok/auth/model"
	"github.com/Thevtok/auth/service"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type LoginJwt struct {
	service service.LoginService
	jwtKey  []byte
}

var jwtKey = []byte("thevtok")

func generateToken(user *model.User) (string, error) {
	// Set token claims
	claims := jwt.MapClaims{}
	claims["username"] = user.Username
	claims["c_username"] = user.C_Username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token expires in 24 hours

	// Create token with claims and secret key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func AuthMiddleware(jwtKey []byte) gin.HandlerFunc {
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
		username := claims["username"].(string)
		c_username := claims["c_username"].(string)

		c.Set("username", username)
		c.Set("c_username", c_username)

		c.Next()
	}
}

func (lj *LoginJwt) Login(c *gin.Context) {
	var credentials model.User

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	user, err := lj.service.LoginSkuy(credentials.Username, credentials.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := generateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func NewUserJwt(s service.LoginService) *LoginJwt {
	loginj := LoginJwt{
		service: s,
		jwtKey:  jwtKey,
	}
	return &loginj
}
