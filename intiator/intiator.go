package intiator

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"restaurant/internal/constant/model/persistencedb"
	"restaurant/internal/handler/middleware"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func Intiate() {
	sampleLogger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Printf(`{"level":"fatal","msg":"failed to initialize sample logger: %v"}
`, err)
		os.Exit(1)
	}

	sampleLogger.Info("initializing config")
	// configName := "config"
	// if name := os.Getenv("CONFIG_NAME"); name != "" {
	// 	configName = name
	// 	sampleLogger.Info(fmt.Sprintf("config name is set to %s", configName))
	// } else {
	// 	sampleLogger.Info("using default config name 'config'")
	// }

	_ = godotenv.Load("./config/.env")

	// InitConfig(configName, "config", sampleLogger)
	sampleLogger.Info("config initialized")

	sampleLogger.Info("initializing database")
	pgxConn := InitDB(os.Getenv("DATABASE_URL"), sampleLogger)
	sampleLogger.Info("database initialized")

	sampleLogger.Info("intializing cache")
	cache := InitCache(os.Getenv("REDIS_URL"), sampleLogger)
	sampleLogger.Info("cache intialized")

	// if viper.GetBool("migration.active")
	if strings.ToLower(os.Getenv("MIGRATION_ACTIVE")) == "true" {
		sampleLogger.Info("initializing migration")
		m := InitiateMigration(os.Getenv("MIGRATION_PATH"), os.Getenv("DATABASE_URL"), sampleLogger)
		UpMigration(m, sampleLogger)
		sampleLogger.Info("migration initialized")
	}

	sampleLogger.Info("initializing persistence layer")
	persistence := InitPersistence(persistencedb.New(pgxConn, sampleLogger), sampleLogger)
	sampleLogger.Info("persistence layer initialized")

	sampleLogger.Info("initializing cache layer")
	cacheLayer := InitCacheLayer(CacheOptions{Redis: cache,
		OrderExpirationTime: 5 * time.Minute}, sampleLogger)
	sampleLogger.Info("cache layer initialized")

	sampleLogger.Info("initializing module")
	module := InitModule(persistence, cacheLayer, sampleLogger)
	sampleLogger.Info("module initialized")

	sampleLogger.Info("initializing handler")
	serverTimeout, err := time.ParseDuration(os.Getenv("SERVER_TIMEOUT"))
	if err != nil {
		sampleLogger.Fatal("unable to parse server timeout duration", zap.Error(err))
	}
	handler := InitHandler(module, sampleLogger, serverTimeout)
	sampleLogger.Info("handler initialized")

	sampleLogger.Info("initializing server")
	server := gin.New()
	gin.SetMode(gin.DebugMode)
	server.Use(middleware.ErrorHandler())
	sampleLogger.Info("server initialized")

	sampleLogger.Info("initializing router")
	v1 := server.Group("/v1")
	InitRouter(v1, handler, module, sampleLogger)
	sampleLogger.Info("router initialized")

	readHeaderTimeout, err := time.ParseDuration(os.Getenv("SERVER_READ_HEADER_TIMEOUT"))
	if err != nil {
		sampleLogger.Fatal("unable to parse read header timeout duration", zap.Error(err))
	}
	srv := &http.Server{
		Addr:              os.Getenv("SERVER_HOST") + ":" + os.Getenv("SERVER_PORT"),
		ReadHeaderTimeout: readHeaderTimeout,
		Handler:           server,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	serverPort, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil {
		sampleLogger.Fatal("unable to parse server port", zap.Error(err))
	}
	sampleLogger.Info("server started",
		zap.String("host", os.Getenv("SERVER_HOST")),
		zap.Int("port", serverPort),
		zap.Time("start_time", time.Now()))

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			sampleLogger.Fatal("error listening to the server", zap.Error(err))
		}
	}()

	// Wait for termination signal
	sig := <-quit
	sampleLogger.Info("server shutting down", zap.String("signal", sig.String()))

	// timeout := viper.GetDuration("server.timeout")
	if serverTimeout == 0 {
		serverTimeout = 5 * time.Second // Default to 5 seconds
		sampleLogger.Warn("server timeout not set, using default", zap.Duration("timeout", serverTimeout))
	}
	ctx, cancel := context.WithTimeout(context.Background(), serverTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		sampleLogger.Fatal("error while shutting down server", zap.Error(err))
	}

	sampleLogger.Info("server shutdown complete")

}
