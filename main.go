package main

import (
	"backend/args_parser"
	"backend/config"
	"backend/logger"
	"backend/src/handlers"
	"backend/src/middleware"
	"backend/src/models"
	"backend/src/repo"
	"backend/src/services"
	"backend/src/utils"
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/stdlib"
)

func main() {
	debugMode := true

	args, err := args_parser.Parse(os.Args)
	if err != nil {
		panic(err)
	}

	logger, err := logger.New(logger.NewLoggerOpts{
		Debug:      debugMode,
		OutputFile: "./log.txt",
	})
	if err != nil {
		panic(err)
	}

	conf, err := config.NewFromFile(args.GetConfigPath())
	if err != nil {
		logger.Fatal().Err(err).Msg("failed parsing config file")
	}

	var key *rsa.PrivateKey
	{
		keyRawBytes, err := os.ReadFile(conf.GetJwtSigningKey())
		if err != nil {
			logger.Fatal().Err(err).Msg("failed reading signing key file")
		}

		keyPem, _ := pem.Decode(keyRawBytes)
		key, err = x509.ParsePKCS1PrivateKey(keyPem.Bytes)
		if err != nil {
			logger.Fatal().Err(err).Msg("failed parsing signing key")
		}
	}

	var sqlDb *sql.DB
	{
		pgConnStr := conf.GetPostgresUrl()
		connConf, err := pgx.ParseConnectionString(pgConnStr)
		if err != nil {
			logger.Fatal().Err(err).Msg("failed parsing postgres connection string")
		}

		sqlDb := stdlib.OpenDB(connConf)
		if err := sqlDb.Ping(); err != nil {
			logger.Fatal().Err(err).Msg("failed pinging postgres db")
		}
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

	if !debugMode {
		gin.SetMode(gin.ReleaseMode)
	}

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
