package main

import (
	_ "WebAnalyzer/docs"
	"WebAnalyzer/internal/config"
	"WebAnalyzer/internal/handler"
	"WebAnalyzer/internal/migration"
	"WebAnalyzer/internal/repository"
	"WebAnalyzer/internal/service"
	"WebAnalyzer/migrations"
	"context"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title WebAnalyzer API
// @version 1.0
// @description This is a sample server for analyzing webpages.
// @BasePath /

var logger *zap.SugaredLogger

func main() {
	logg, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize zap logger: %v", err))
	}
	defer logg.Sync()
	logger = logg.Sugar()

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	conf := mustLoadConfig(getConfigPath())
	mustMigratePostgres(conf)
	db := mustInitPostgres(conf)
	repos := repository.NewRepositoryContainer(db)
	services := service.NewServiceContainer(repos)

	httpHandler := handler.NewHttpHandler(services, conf, logger)
	httpHandler.Init()

	setupGracefulShutdown(ctx, cancelFunc, httpHandler)
}

func mustLoadConfig(path string) *config.Config {
	conf, err := config.Load(path)
	if err != nil {
		logger.Fatalf("Cannot read config: %v", err)
	}

	dbUrl := os.Getenv("DATABASE_URL")

	if dbUrl != "" {
		conf.DB.URL = dbUrl
	}

	return conf
}

func getConfigPath() string {
	if path := os.Getenv("CONFIG_PATH"); path != "" {
		return path
	}
	return "config"
}

func mustMigratePostgres(conf *config.Config) {
	dbFiles := migrations.GetPostgresMigrations()
	dbVersion, err := migration.PostgresMigrate(conf.DB.URL, conf.Migration, dbFiles)
	if err != nil {
		logger.Fatalf("Cannot migrate db: %v", err)
	}
	logger.Infof("dbVersion: %v", dbVersion)
}

func mustInitPostgres(conf *config.Config) *gorm.DB {
	db, err := repository.InitORM(conf.DB)
	if err != nil {
		logger.Fatalf("Cannot init db: %v", err)
	}
	return db
}

func setupGracefulShutdown(ctx context.Context, cancelFunc context.CancelFunc, handler handler.HttpHandler) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer signal.Stop(sigChan)

	go func() {
		<-sigChan
		logger.Info("Signal received. Initiating shutdown...")
		cancelFunc()
	}()

	<-ctx.Done()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	logger.Info("Shutting down HTTP server...")
	if err := handler.Stop(shutdownCtx); err != nil {
		logger.Fatalf("HTTP server shutdown error: %v", err)
	}

	logger.Info("Shutdown complete.")
}
