package src

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
	"log"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
	"time"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/stdlib"
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

	//-----------------------------------------

	args, err := args_parser.Parse(osArgs)
	if err != nil {
		log.Fatalf("failed to parse os args: %v\n", err)
	}

	logger, err := logger.New(logger.NewLoggerOpts{
		Debug:      debugMode,
		OutputFile: args.GetLogPath(),
	})
	if err != nil {
		log.Fatalf("failed to create logger object: %v\n", err)
	}

	conf, err := config.NewFromFile(args.GetConfigPath())
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to parse config file")
	}

	//-----------------------------------------

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

	ctx, stop := signal.NotifyContext(ctx, signals...)
	defer stop()

	var sqlDb *sql.DB // TODO: move to integrations package
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

	// Build business-logic objects
	var (
		userService      services.UserService
		shortlinkService services.ShortlinkService
	)
	{
		var (
			jwtUtil      = utils.NewJwtUtil(key)
			passwordUtil = utils.NewPasswordUtil()

			userRepo        = repos.NewUserRepo(sqlDb)
			emailRepo       = repos.NewEmailRepo()
			actionTokenRepo = repos.NewActionTokenRepo(sqlDb)
			linksCache      = repos.NewCacheInmem[string, string](7 * 24 * 60 * 60)
			userCache       = repos.NewCacheInmem[string, models.UserDTO](60 * 60)
		)

		// Periodically trigger cache cleanup
		go func() {
			tmr := time.NewTicker(5 * time.Minute)
			defer tmr.Stop()

			for {
				select {
				case <-ctx.Done():
					return
				case <-tmr.C:
					userCache.CheckExpired()
					linksCache.CheckExpired()
				}
			}
		}()

		userService = services.NewUserService(
			services.UserServiceDeps{
				Jwt:             jwtUtil,
				Password:        passwordUtil,
				UserRepo:        userRepo,
				UserCache:       userCache,
				EmailRepo:       emailRepo,
				ActionTokenRepo: actionTokenRepo,
			},
		)
		shortlinkService = services.NewShortlinkSevice(
			services.NewShortlinkServiceParams{
				Cache: linksCache,
			},
		)
	}

	clientNotifier := client_notifier.NewBasicNotifier()

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

	srv := server.New(
		server.NewServerOpts{
			DebugMode:        debugMode,
			Logger:           logger,
			Notifier:         clientNotifier,
			ShortlinkService: shortlinkService,
			UserService:      userService,
		},
	)
	srv.Run(ctx, conf.GetPort())
}
