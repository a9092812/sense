package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/Kartik30R/sense/config"
	"github.com/Kartik30R/sense/internal/handlers"
	"github.com/Kartik30R/sense/internal/kafka"
	"github.com/Kartik30R/sense/internal/models"
	"github.com/Kartik30R/sense/internal/pipeline"
	"github.com/Kartik30R/sense/internal/redis"
	"github.com/Kartik30R/sense/internal/repository"
	routes "github.com/Kartik30R/sense/internal/router"
	services "github.com/Kartik30R/sense/internal/service"
	"github.com/Kartik30R/sense/internal/websocket"
	"github.com/Kartik30R/sense/pkg/logger"
)

func main() {

	// 1️⃣ Initialize Logger
	logger.InitLogger()
	defer logger.Log.Sync()

	// 2️⃣ Load Configuration
	cfg := config.LoadAllConfig()
	logger.Info("Starting Sense Platform", zap.String("env", cfg.App.Env))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// 3️⃣ Infrastructure
	config.InitDB(cfg.DB)
	config.InitRedis(cfg.Redis)

	// Run migrations
	if err := models.RunMigrations(config.DB); err != nil {
		logger.Fatal("Database migration failed", zap.Error(err))
	}

	// 4️⃣ Repository Layer
	repos := repository.NewRegistry(config.DB)

	// 5️⃣ Pipeline
	jobChannel := pipeline.NewJobChannel(10000)
	batchAggregator := pipeline.NewBatchAggregator(repos)
	redisPub := redis.NewPublisher()

	workerPool := pipeline.NewWorkerPool(
		cfg.App.WorkerCount,
		jobChannel,
		batchAggregator,
		repos,
	)

	workerPool.Start(ctx)

	// 6️⃣ Kafka
	if err := kafka.EnsureTopics(cfg.Kafka.Brokers[0], kafka.SensorPacketsTopic, 1); err != nil {
		logger.Warn("Failed to ensure Kafka topics", zap.Error(err))
	}

	producer := kafka.NewProducer(cfg.Kafka)
	defer producer.Close()

	consumer := kafka.NewConsumer(cfg.Kafka, cfg.RateLimit)
	defer consumer.Close()

	// Start Kafka consumer
	go consumer.Start(ctx, jobChannel.SendChannel())

	// 7️⃣ WebSocket Hubs

	// Dashboard hub
	hub := websocket.NewHub()
	go hub.Run()

	// Redis → WebSocket bridge
	sub := redis.NewSubscriber()
	go func() {
		ch := sub.Subscribe(ctx, redis.SensorDataChannel)
		for msg := range ch {
			if raw, ok := msg.([]byte); ok {
				hub.Broadcast(raw)
			}
		}
	}()

	// 8️⃣ Services

	deviceService := services.NewDeviceService(repos.Device)
	mobileDeviceService := services.NewMobileDeviceService(repos.MobileDevice)

	packetService := services.NewPacketService(
		producer,
		deviceService,
		mobileDeviceService,
	)
	historyService := services.NewHistoryService(config.DB)

	// Ingest hub (mobile live stream)
	ingestHub := websocket.NewIngestHub(redisPub, packetService, mobileDeviceService)
	go ingestHub.Run()

	// 9️⃣ Handlers

	packetHandler := handlers.NewPacketHandler(packetService)
	historyHandler := handlers.NewHistoryHandler(historyService)
	deviceHandler := handlers.NewDeviceHandler(deviceService)
	mobileDeviceHandler := handlers.NewMobileDeviceHandler(mobileDeviceService, deviceService)

	// 🔟 HTTP Server

	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"*"},
	}))

	routes.SetupRoutes(
		router,
		packetHandler,
		historyHandler,
		deviceHandler, // NEW
		mobileDeviceHandler,
		hub,
		ingestHub,
	)

	srvAddr := ":" + cfg.App.Port

	logger.Info("HTTP server starting", zap.String("addr", srvAddr))

	go func() {
		if err := router.Run("0.0.0.0:" + cfg.App.Port); err != nil {
			logger.Fatal("HTTP server failed", zap.Error(err))
		}
	}()

	<-ctx.Done()

	logger.Info("Shutting down gracefully...")
}
 