package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	App       AppConfig       `mapstructure:"app"`
	Database  DatabaseConfig  `mapstructure:"database"`
	Scanner   ScannerConfig   `mapstructure:"scanner"`
	Scraper   ScraperConfig   `mapstructure:"scraper"`
	Stream    StreamConfig    `mapstructure:"stream"`
	Log       LogConfig       `mapstructure:"log"`
	CORS      CORSConfig      `mapstructure:"cors"`
	RateLimit RateLimitConfig `mapstructure:"rate_limit"`
}

type AppConfig struct {
	Name            string `mapstructure:"name"`
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	Mode            string `mapstructure:"mode"`
	ContextTimeout  int    `mapstructure:"context_timeout"`
}

type DatabaseConfig struct {
	Path             string `mapstructure:"path"`
	MaxOpenConns     int    `mapstructure:"max_open_conns"`
	MaxIdleConns     int    `mapstructure:"max_idle_conns"`
	ConnMaxLifetime  int    `mapstructure:"conn_max_lifetime"`
}

type ScannerConfig struct {
	MaxDepth        int      `mapstructure:"max_depth"`
	ExcludePatterns []string `mapstructure:"exclude_patterns"`
	FileSizeLimit   int64    `mapstructure:"file_size_limit"`
}

type ScraperConfig struct {
	TMDBAPIKey     string `mapstructure:"tmdb_api_key"`
	TMDBBaseURL    string `mapstructure:"tmdb_base_url"`
	Language       string `mapstructure:"language"`
	Timeout        int    `mapstructure:"timeout"`
	RetryCount     int    `mapstructure:"retry_count"`
}

type StreamConfig struct {
	HLSDir              string `mapstructure:"hls_dir"`
	DefaultQuality      string `mapstructure:"default_quality"`
	TranscodeEnabled    bool   `mapstructure:"transcode_enabled"`
	HardwareAcceleration bool  `mapstructure:"hardware_acceleration"`
	MaxConcurrent       int    `mapstructure:"max_concurrent"`
}

type LogConfig struct {
	Level     string `mapstructure:"level"`
	Path      string `mapstructure:"path"`
	MaxSize   int    `mapstructure:"max_size"`
	MaxAge    int    `mapstructure:"max_age"`
	MaxBackups int   `mapstructure:"max_backups"`
}

type CORSConfig struct {
	AllowOrigins []string `mapstructure:"allow_origins"`
	AllowMethods []string `mapstructure:"allow_methods"`
	AllowHeaders []string `mapstructure:"allow_headers"`
}

type RateLimitConfig struct {
	RPS    float64 `mapstructure:"rps"`
	Bursts int     `mapstructure:"bursts"`
}

func LoadConfig() (*Config, error) {
	// 加载.env文件
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: .env file not found")
	}

	// 设置默认配置文件路径
	configPath := "./config/config.yaml"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// 尝试从上级目录加载
		configPath = "./backend/config/config.yaml"
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			// 尝试从当前目录加载
			configPath = "config.yaml"
			if _, err := os.Stat(configPath); os.IsNotExist(err) {
				return nil, fmt.Errorf("config file not found")
			}
		}
	}

	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	// 自动环境变量
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// 解析配置
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// 应用环境变量覆盖
	if port := os.Getenv("API_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.App.Port = p
		}
	}

	if dbPath := os.Getenv("DATABASE_PATH"); dbPath != "" {
		cfg.Database.Path = dbPath
	}

	if tmdbAPIKey := os.Getenv("TMDB_API_KEY"); tmdbAPIKey != "" {
		cfg.Scraper.TMDBAPIKey = tmdbAPIKey
	}

	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		cfg.Log.Level = logLevel
	}

	// 确保目录存在
	ensureDirectories(cfg)

	return &cfg, nil
}

func ensureDirectories(cfg Config) {
	directories := []string{
		filepath.Dir(cfg.Database.Path),
		cfg.Log.Path,
		cfg.Stream.HLSDir,
		"./data/posters",
		"./data/backdrops",
	}

	for _, dir := range directories {
		if dir != "" {
			os.MkdirAll(dir, 0755)
		}
	}
}
