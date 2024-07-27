package main

import (
	"backend/args_parser"
	"backend/config"
	"backend/src/handlers"
	"backend/src/middleware"
	"backend/src/models"
	"backend/src/repo"
	"backend/src/services"
	"backend/src/utils"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/stdlib"
)

func main() {
	args, err := args_parser.Parse(os.Args)
	if err != nil {
		panic(err)
	}

	conf, err := config.NewFromFile(args.GetConfigPath())
	if err != nil {
		panic(err)
	}

	keyRawBytes, err := os.ReadFile(conf.GetJwtSigningKey())
	if err != nil {
		panic(err)
	}

	keyPem, _ := pem.Decode(keyRawBytes)
	key, err := x509.ParsePKCS1PrivateKey(keyPem.Bytes)
	if err != nil {
		panic(err)
	}

	pgConnStr := conf.GetPostgresUrl()
	connConf, err := pgx.ParseConnectionString(pgConnStr)
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

	r.Run(fmt.Sprintf(":%d", conf.GetPort()))
}
