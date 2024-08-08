package main

import (
	"backend/src/args_parser"
	"backend/src/config"
	"backend/src/core/models"
	"backend/src/core/repos"
	"backend/src/core/services"
	"backend/src/core/utils"
	"backend/src/logger"
	"backend/src/server/handlers"
	"backend/src/server/middleware"
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
		OutputFile: args.GetLogPath(),
	})
	if err != nil {
		panic(err)
	}

	logger.Log().Msg("initializing service...")

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

		sqlDb = stdlib.OpenDB(connConf)
		if err := sqlDb.Ping(); err != nil {
			logger.Fatal().Err(err).Msg("failed pinging postgres db")
		}
	}

	jwtUtil := utils.NewJwtUtil(key)
	passwordUtil := utils.NewPasswordUtil()
	userRepo := repos.NewUserRepo(sqlDb)
	userCache := repos.NewCacheInmem[string, models.UserDTO](60 * 60)
	emailRepo := repos.NewEmailRepo()
	actionTokenRepo := repos.NewActionTokenRepo(sqlDb)

	userService := services.NewUserService(
		services.UserServiceDeps{
			Jwt:             jwtUtil,
			Password:        passwordUtil,
			UserRepo:        userRepo,
			UserCache:       userCache,
			EmailRepo:       emailRepo,
			ActionTokenRepo: actionTokenRepo,
		},
	)
	linkService := services.NewShortlinkSevice(
		services.NewShortlinkServiceParams{
			Cache: repos.NewCacheInmem[string, string](7 * 24 * 60 * 60),
		},
	)

	if !debugMode {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(middleware.NewRequestLogMiddleware(logger))
	r.Use(gin.Recovery())

	linkGroup := r.Group("/s")
	linkGroup.POST("/new", handlers.NewShortlinkCreateHandler(linkService))
	linkGroup.GET("/:linkId", handlers.NewShortlinkResolveHandler(linkService))

	userGroup := r.Group("/user")
	userGroup.POST("/create", handlers.NewUserCreateHandler(userService))
	userGroup.POST("/login", handlers.NewUserLoginHandler(userService))

	dummyGroup := r.Group("/dummy")
	dummyGroup.Use(middleware.NewAuthMiddleware(userService))
	dummyGroup.GET("/", handlers.NewDummyHandler())

	listenAddr := fmt.Sprintf(":%d", conf.GetPort())
	logger.Log().Msgf("server listening on %s", listenAddr)

	err = r.Run(listenAddr)
	if err != nil {
		logger.Fatal().Err(err).Msg("server stopped with error")
	}
}
