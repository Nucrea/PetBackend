package main

import (
	httpserver "backend/internal/http_server"
	"backend/internal/integrations"
	"backend/pkg/logger"
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
)

func main() {
	ctx := context.Background()

	config, err := LoadConfig("config.yaml")
	if err != nil {
		log.Fatal(err.Error())
	}

	logger, err := logger.New(ctx, logger.NewLoggerOpts{
		Debug:      true,
		OutputFile: config.App.LogFile,
	})
	if err != nil {
		logger.Fatal().Err(err)
	}

	emailer, err := NewEmailer(config.SMTP)
	if err != nil {
		logger.Fatal().Err(err)
	}

	metrics := integrations.NewMetrics("notifyer")

	ginRouter := gin.New()
	ginRouter.GET("/metrics", gin.WrapH(metrics.HttpHandler()))
	ginRouter.GET("/health", func(ctx *gin.Context) {
		ctx.Status(200)
	})

	kafkaReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: config.Kafka.Brokers,
		Topic:   config.Kafka.Topic,
		GroupID: config.Kafka.ConsumerGroupId,
	})
	kafkaReader.SetOffset(kafka.LastOffset)

	eventHandler := NewEventHandler(config, logger, metrics, emailer)
	go eventHandler.eventLoop(ctx, kafkaReader)

	logger.Log().Msg("notifyer service started")

	srv := httpserver.New(
		httpserver.NewServerOpts{
			Logger:     logger,
			HttpServer: ginRouter,
		},
	)
	srv.Run(ctx, config.App.Port)
}
