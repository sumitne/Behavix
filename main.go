package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"behavix-ai/internal/app"
	"behavix-ai/internal/config"
	httpeventsdebug "behavix-ai/internal/infrastructure/http/eventsdebug"
	httpingestion "behavix-ai/internal/infrastructure/http/ingestion"
	httpinsight "behavix-ai/internal/infrastructure/http/insight"
	httpserver "behavix-ai/internal/infrastructure/http/server"
	"behavix-ai/internal/infrastructure/postgres"
	"behavix-ai/internal/logger"
	serviceingestion "behavix-ai/internal/service/ingestion"
	serviceinsight "behavix-ai/internal/service/insight"

	"go.uber.org/zap"
)

func main() {
	if len(os.Args) >= 2 && os.Args[1] == "migrate-up" {
		app.RunMigrate("migrations")
		return
	}

	runServer()
}

func runServer() {
	ctx := context.Background()
	log, err := logger.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("config load failed", zap.Error(err))
	}

	pool, err := postgres.New(ctx, cfg.PostgresDSN)
	if err != nil {
		log.Fatal("postgres connect failed", zap.Error(err))
	}
	defer pool.Close()

	tenantRepo := postgres.NewTenantRepository(pool)
	eventRepo := postgres.NewEventRepository(pool)
	insightRepo := postgres.NewInsightRepository(pool)
	ingestionSvc := serviceingestion.NewService(eventRepo)
	insightSvc := serviceinsight.NewService(insightRepo)
	ingestionHandler := httpingestion.NewHandler(ingestionSvc)
	insightHandler := httpinsight.NewHandler(insightSvc)
	eventsDebugHandler := httpeventsdebug.NewHandler(eventRepo)

	router := httpserver.NewRouter(httpserver.RouterConfig{
		Logger:         log,
		EventIngestion: ingestionHandler,
		EventList:      eventsDebugHandler,
		InsightList:    insightHandler,
		TenantRepo:     tenantRepo,
		UseTenantAuth:  true,
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: router,
	}

	go func() {
		log.Info("server listening", zap.Int("port", cfg.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("server shutdown failed", zap.Error(err))
	}
	log.Info("server stopped")
}
