package main

import (
	"backend/internal/core/repos"
	"backend/internal/core/services"
	grpcserver "backend/internal/grpc_server"
	"backend/internal/grpc_server/shortlinks"
	httpserver "backend/internal/http_server"
	"backend/internal/http_server/middleware"
	"backend/internal/integrations"
	"backend/pkg/cache"
	"backend/pkg/logger"
	"context"
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

var rootCmd = &cobra.Command{
	Use:   "shortlinks",
	Short: "shortlinks is a microservice for creating ang managing shortlinks",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func main() {
	ctx := context.Background()

	var (
		configPath = ""
		logPath    = ""
	)
	{
		rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "path to configuration file")
		rootCmd.MarkPersistentFlagRequired("config")

		rootCmd.PersistentFlags().StringVarP(&logPath, "logfile", "l", "", "path to log file")
		rootCmd.MarkPersistentFlagRequired("logfile")

		if err := rootCmd.Execute(); err != nil {
			panic(err)
		}
	}

	log, err := logger.New(ctx, logger.NewLoggerOpts{
		Debug:      true,
		OutputFile: logPath,
	})
	if err != nil {
		panic(err)
	}

	conf, err := LoadConfig(configPath)
	if err != nil {
		log.Error().Err(err).Msg("failed loading config")
	}

	pgDb, err := integrations.NewPostgresConn(ctx, conf.GetPostgresUrl())
	if err != nil {
		log.Error().Err(err).Msg("failed connecting to postgres")
	}

	tracer, err := integrations.NewTracer("backend")
	if err != nil {
		log.Fatal().Err(err).Msg("failed initializing tracer")
	}

	repo := repos.NewShortlinkRepo(pgDb, tracer)
	service := services.NewShortlinkSevice(
		services.NewShortlinkServiceParams{
			Cache: cache.NewCacheInmem[string, string](),
			Repo:  repo,
		},
	)

	RunServer(ctx, log, tracer, conf, service)
}

func RunServer(ctx context.Context, log logger.Logger, tracer trace.Tracer, conf IConfig, shortlinkService services.ShortlinkService) {
	host := fmt.Sprintf("http://localhost:%d", conf.GetHttpPort())
	debugMode := true
	if !debugMode {
		gin.SetMode(gin.ReleaseMode)
	}

	prometheus := integrations.NewPrometheus()

	r := gin.New()
	r.Any("/metrics", gin.WrapH(prometheus.GetRequestHandler()))
	r.GET("/health", func(ctx *gin.Context) {
		ctx.Status(200)
	})

	r.Use(middleware.NewRecoveryMiddleware(log, prometheus, debugMode))
	r.Use(middleware.NewRequestLogMiddleware(log, tracer, prometheus))
	r.Use(middleware.NewTracingMiddleware(tracer))

	linkGroup := r.Group("/s")
	linkGroup.POST("/new", NewShortlinkCreateHandler(log, shortlinkService, host))
	linkGroup.GET("/:linkId", NewShortlinkResolveHandler(log, shortlinkService))

	grpcObj := &ShortlinksGrpc{
		log:              log,
		host:             host,
		shortlinkService: shortlinkService,
	}

	grpcUnderlying := grpc.NewServer()
	shortlinks.RegisterShortlinksServer(grpcUnderlying, grpcObj)

	httpServer := httpserver.New(
		httpserver.NewServerOpts{
			Logger:     log,
			HttpServer: r,
		},
	)
	grpcServer := grpcserver.New(
		grpcserver.NewServerOpts{
			Logger:     log,
			GrpcServer: grpcUnderlying,
		},
	)

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		httpServer.Run(ctx, conf.GetHttpPort())
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		grpcServer.Run(ctx, conf.GetGrpcPort())
	}()

	wg.Wait()
}

// func RunTestGrpcClient() {
// 	go func() {
// 		conn, err := grpc.NewClient(
// 			fmt.Sprintf(":%d", conf.GetPort()),
// 			grpc.WithTransportCredentials(insecure.NewCredentials()),
// 		)
// 		if err != nil {
// 			log.Fatal().Err(err).Msg("failed initializing grpc test client")
// 		}
// 		defer conn.Close()

// 		c := shortlinks.NewShortlinksClient(conn)

// 		for {
// 			select {
// 			case <-ctx.Done():
// 				return
// 			default:
// 			}

// 			res, err := c.Create(ctx, &shortlinks.CreateRequest{
// 				Url: "https://google.com",
// 			})
// 			if err != nil {
// 				log.Error().Err(err).Msg("failed creating shortlink")
// 			} else {
// 				log.Log().Msgf("Successfully created link: %s", res.GetLink())
// 			}

// 			time.Sleep(3 * time.Second)
// 		}
// 	}()
// }
