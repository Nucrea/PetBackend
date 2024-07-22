package main

import (
	"backend/src/handlers"
	"backend/src/middleware"
	"backend/src/models"
	"backend/src/repo"
	"backend/src/services"
	"backend/src/utils"
	"crypto/rand"
	"crypto/rsa"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/stdlib"
)

func main() {
	key, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		panic(err)
	}
	// keyBytes, err := x509.MarshalPKCS8PrivateKey(key)
	// if err != nil {
	// 	panic(err)
	// }

	connConf, err := pgx.ParseConnectionString("postgres://postgres:postgres@localhost:5432/postgres")
	if err != nil {
		panic(err)
	}

	sqlDb := stdlib.OpenDB(connConf)
	if err := sqlDb.Ping(); err != nil {
		panic(err)
	}

	jwtUtil := utils.NewJwtUtil(key)
	passwordUtil := utils.NewPasswordUtil()
	userRepo := repo.NewUserRepo(sqlDb)
	userCache := repo.NewCacheInmem[string, models.UserDTO](60 * 60)

	userService := services.NewUserService(
		services.UserServiceDeps{
			Jwt:       jwtUtil,
			Password:  passwordUtil,
			UserRepo:  userRepo,
			UserCache: userCache,
		},
	)

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	userGroup := r.Group("/user")
	userGroup.POST("/create", handlers.NewUserCreateHandler(userService))
	userGroup.POST("/login", handlers.NewUserLoginHandler(userService))

	dummyGroup := r.Group("/dummy")
	dummyGroup.Use(middleware.NewAuthMiddleware(userService))
	dummyGroup.GET("/", handlers.NewDummyHandler())

	r.Run(":8080")
}
