package main

import (
	"backend/cmd/backend/server"
	"backend/internal/core/models"
	"backend/internal/core/repos"
	"backend/internal/core/services"
	"backend/internal/core/utils"
	"backend/internal/integrations"
	"backend/pkg/cache"
	"backend/pkg/logger"
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
	"time"
)

type App struct{}

type RunParams struct {
	Ctx     context.Context
	OsArgs  []string
	EnvVars map[string]string
}

func (a *App) Run(p RunParams) {
	var (
		ctx       = p.Ctx
		osArgs    = p.OsArgs
		_         = p.EnvVars
		debugMode = false // TODO: replace with flag from conf
	)

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

	ctx, stop := signal.NotifyContext(ctx, signals...)
	defer stop()

	//-----------------------------------------

	args, err := CmdArgsParse(osArgs)
	if err != nil {
		log.Fatalf("failed to parse os args: %v\n", err)
	}

	logger, err := logger.New(
		ctx,
		logger.NewLoggerOpts{
			Debug:      debugMode,
			OutputFile: args.GetLogPath(),
		},
	)
	if err != nil {
		log.Fatalf("failed to create logger object: %v\n", err)
	}

	conf, err := LoadConfig(args.GetConfigPath())
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to parse config file")
	}

	//-----------------------------------------

	logger.Log().Msg("initializing service...")
	defer logger.Log().Msg("service stopped")

	sqlDb, err := integrations.NewPostgresConn(ctx, conf.GetPostgresUrl())
	if err != nil {
		logger.Fatal().Err(err).Msg("failed connecting to postgres")
	}

	kafka := integrations.NewKafka(conf.GetKafkaUrl(), conf.GetKafkaTopic())

	var key *rsa.PrivateKey
	{
		keyRawBytes, err := os.ReadFile(args.GetSigningKeyPath())
		if err != nil {
			logger.Fatal().Err(err).Msg("failed reading signing key file")
		}

		keyPem, _ := pem.Decode(keyRawBytes)
		key, err = x509.ParsePKCS1PrivateKey(keyPem.Bytes)
		if err != nil {
			logger.Fatal().Err(err).Msg("failed parsing signing key")
		}
	}

	tracer, err := integrations.NewTracer("backend")
	if err != nil {
		logger.Fatal().Err(err).Msg("failed initializing tracer")
	}

	// Build business-logic objects
	var (
		userService      services.UserService
		shortlinkService services.ShortlinkService
	)
	{
		var (
			jwtUtil      = utils.NewJwtUtil(key)
			passwordUtil = utils.NewPasswordUtil()

			userRepo        = repos.NewUserRepo(sqlDb, tracer)
			actionTokenRepo = repos.NewActionTokenRepo(sqlDb)
			shortlinkRepo   = repos.NewShortlinkRepo(sqlDb, tracer)
			eventRepo       = repos.NewEventRepo(kafka)

			userCache          = cache.NewCacheInmemSharded[models.UserDTO](cache.ShardingTypeInteger)
			loginAttemptsCache = cache.NewCacheInmem[string, int]()
			jwtCache           = cache.NewCacheInmemSharded[string](cache.ShardingTypeJWT)
			linksCache         = cache.NewCacheInmem[string, string]()
		)

		// Periodically trigger cache cleanup
		go func() {
			tmr := time.NewTicker(15 * time.Minute)
			defer tmr.Stop()

			for {
				select {
				case <-ctx.Done():
					return
				case <-tmr.C:
					userCache.CheckExpired()
					jwtCache.CheckExpired()
					linksCache.CheckExpired()
					loginAttemptsCache.CheckExpired()
				}
			}
		}()

		userService = services.NewUserService(
			services.UserServiceDeps{
				Jwt:                jwtUtil,
				Password:           passwordUtil,
				UserRepo:           userRepo,
				UserCache:          userCache,
				LoginAttemptsCache: loginAttemptsCache,
				JwtCache:           jwtCache,
				EventRepo:          *eventRepo,
				ActionTokenRepo:    actionTokenRepo,
				Logger:             logger,
			},
		)
		shortlinkService = services.NewShortlinkSevice(
			services.NewShortlinkServiceParams{
				Cache: linksCache,
				Repo:  shortlinkRepo,
			},
		)

		// TODO: Run cleanup routine
		// go shortlinkService.ShortlinkRoutine(ctx)
	}

	// Start profiling
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

	srv := server.NewServer(
		server.NewServerOpts{
			DebugMode:        debugMode,
			Logger:           logger,
			ShortlinkService: shortlinkService,
			UserService:      userService,
			Tracer:           tracer,
		},
	)
	srv.Run(ctx, conf.GetPort())
}
