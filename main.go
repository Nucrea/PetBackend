package main

import (
	"backend/src/args_parser"
	"backend/src/client_notifier"
	"backend/src/config"
	"backend/src/core/models"
	"backend/src/core/repos"
	"backend/src/core/services"
	"backend/src/core/utils"
	"backend/src/logger"
	"backend/src/server"
	"context"
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
	"time"

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
	defer logger.Log().Msg("service stopped")

	signals := []os.Signal{
		os.Kill,
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGKILL,
		syscall.SIGSTOP,
		syscall.SIGQUIT,
		syscall.SIGABRT,
		syscall.SIGHUP,
	}

	ctx, stop := signal.NotifyContext(context.Background(), signals...)
	defer stop()

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
	emailRepo := repos.NewEmailRepo()
	actionTokenRepo := repos.NewActionTokenRepo(sqlDb)

	userCache := repos.NewCacheInmem[string, models.UserDTO](60 * 60)
	go func() {
		tmr := time.NewTicker(30 * time.Second)
		defer tmr.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-tmr.C:
				userCache.CheckExpired()
			}
		}
	}()

	clientNotifier := client_notifier.NewBasicNotifier()

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

	if args.GetProfilePath() != "" {
		pprofFile, err := os.Create(args.GetProfilePath())
		if err != nil {
			logger.Fatal().Err(err).Msg("can not create profile file")
		}

		if err := pprof.StartCPUProfile(pprofFile); err != nil {
			logger.Fatal().Err(err).Msg("can not start cpu profiling")
		}
		defer func() {
			logger.Log().Msg("stopping profiling...")
			pprof.StopCPUProfile()
			pprofFile.Close()
		}()
	}

	srv := server.New(
		server.NewServerOpts{
			DebugMode:        debugMode,
			Logger:           logger,
			Notifier:         clientNotifier,
			ShortlinkService: linkService,
			UserService:      userService,
		},
	)
	srv.Run(ctx, conf.GetPort())
}
