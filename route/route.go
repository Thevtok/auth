package route

import (
	"log"

	"github.com/Thevtok/auth/cont"
	"github.com/Thevtok/auth/db"
	"github.com/Thevtok/auth/repo"
	"github.com/Thevtok/auth/service"
	"github.com/gin-gonic/gin"
)

func Run() {
	database := repo.ConnectDB()
	s_key := []byte(db.DotEnv("SECRET_KEY"))
	authMiddleware := cont.AuthMiddleware(s_key)

	r := gin.Default()
	studentRepo := repo.NewStudentRepo(database)
	loginService := service.NewLoginService(studentRepo)
	loginJwt := cont.NewUserJwt(loginService)

	r.POST("/login", authMiddleware, loginJwt.Login)

	if err := r.Run(db.DotEnv("SERVER_PORT")); err != nil {
		log.Fatal(err)
	}
}
