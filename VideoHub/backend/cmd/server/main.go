package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"videohub/backend/internal/api/router"
	"videohub/backend/internal/pkg/config"
	"videohub/backend/internal/pkg/database"
	"videohub/backend/internal/repository"
	"videohub/backend/internal/service"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化日志
	logger := initLogger(cfg.Log.Level, cfg.Log.Path)
	defer logger.Sync()

	logger.Info("Starting media server")

	// 确保数据目录存在
	ensureDirectories()

	// 连接数据库
	db, err := database.NewSQLite(database.DatabaseConfig{
		Path:             cfg.Database.Path,
		MaxOpenConns:     cfg.Database.MaxOpenConns,
		MaxIdleConns:     cfg.Database.MaxIdleConns,
		ConnMaxLifetime:  cfg.Database.ConnMaxLifetime,
	})
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	// 数据库迁移
	if err := database.Migrate(db); err != nil {
		logger.Fatal("Failed to migrate database", zap.Error(err))
	}

	// 创建索引
	if err := database.CreateIndexes(db); err != nil {
		logger.Warn("Failed to create indexes", zap.Error(err))
	}

	// 初始化仓库
	repo := repository.NewRepository(db)

	// 初始化服务
	svc := service.NewService(
		repo,
		db,
		cfg.Scraper.TMDBAPIKey,
		cfg.Stream.HLSDir,
	)

	// 初始化路由
	r := router.SetupRouter(svc, logger)

	// 启动服务器
	serverAddr := fmt.Sprintf("%s:%d", cfg.App.Host, cfg.App.Port)
	logger.Info("Server starting", zap.String("address", serverAddr))

	if err := r.Run(serverAddr); err != nil {
		logger.Fatal("Server failed", zap.Error(err))
	}
}

func initLogger(level, path string) *zap.Logger {
	var zapLevel zapcore.Level
	switch level {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "info":
		zapLevel = zapcore.InfoLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	default:
		zapLevel = zapcore.InfoLevel
	}

	config := zap.Config{
		Level: zap.NewAtomicLevelAt(zapLevel),
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "message",
			StacktraceKey: "stacktrace",
		},
		OutputPaths:      []string{"stdout", filepath.Join(path, "app.log")},
		ErrorOutputPaths: []string{"stderr", filepath.Join(path, "error.log")},
	}

	logger, _ := config.Build()
	return logger
}

func ensureDirectories() {
	directories := []string{
		"./data",
		"./data/hls",
		"./data/posters",
		"./data/backdrops",
		"./logs",
	}

	for _, dir := range directories {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Printf("Failed to create directory %s: %v", dir, err)
		}
	}
}
