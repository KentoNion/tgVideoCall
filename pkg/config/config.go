package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type DB struct {
	User           string `yaml:"user" env-required:"true"`
	Pass           string `yaml:"password" env-required:"true"`
	Host           string `yaml:"host"`
	Ssl            string `yaml:"sslmode" env-required:"true"`
	MigrationsPath string `yaml:"migrations_path" env-required:"true"`
}

type APIKeys struct {
	Telegram        string `yaml:"telegram" env-required:"true"`
	TelegramPhone   string `yaml:"telegram_phone" env-default:"+79934771502"` // номер для userbot (MTProto), код/2FA — в сообщениях боту
	Youtube         string `yaml:"youtube" env-required:"true"`
	TelegramAppId   int    `yaml:"telegram_app_id" env-required:"true"`
	TelegramAppHash string `yaml:"telegram_app_hash" env-required:"true"`
}

type Log struct {
	FilePath string `yaml:"file_path"`
}

type Service struct {
	StartupLag     time.Duration `yaml:"startup_lag"`
	Cooldown       time.Duration `yaml:"cooldown" default:"3600"`
	Timeout        time.Duration `yaml:"timeout" default:"60"`
	MaxConnections int           `yaml:"max_watcher_connections" default:"100"`
}

type Downloader struct {
	Host          string `yaml:"host" env-required:"true"`
	Max_downloads int    `yaml:"max_parallel_downloads" default:"10"`
	VideoPath     string `yaml:"video_path"`
}

type UploaderConfig struct {
	MaxRetries    int           `env:"UPLOADER_MAX_RETRIES" envDefault:"5"`
	RetryCooldown time.Duration `env:"UPLOADER_RETRY_COOLDOWN" envDefault:"5m"`
}

type Config struct {
	Env        string     `yaml:"env"`
	DB         DB         `yaml:"postgres_db"`
	APIKeys    APIKeys    `yaml:"API_keys"`
	Log        Log        `yaml:"logger"`
	Service    Service    `yaml:"service"`
	Downloader Downloader `yaml:"downloader"`
	Uploader   UploaderConfig
}

func MustLoad() Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "../config.yaml"
	}

	//проверка существует ли файл
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("cannot read config file")
	}

	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatal(err)
	}
	//-----------------------------------------------------------------------------------------------
	dbhost := os.Getenv("DB_HOST") //получаем переменную из окружения (она есть если запущен в докер контейнере)
	if dbhost != "" {
		//time.Sleep(30 * time.Second) //если мы в докер контейнере, дадим время бд чтоб она поднялась
		cfg.DB.Host = dbhost
	}
	downloaderHost := os.Getenv("DOWNLOADER_HOST")
	if downloaderHost != "" {
		cfg.Downloader.Host = downloaderHost
	}
	migrationsPath := os.Getenv("MIGRATIONS_PATH")
	if migrationsPath != "" {
		cfg.DB.MigrationsPath = migrationsPath
	}
	dbUser := os.Getenv("DB_USER")
	if dbUser != "" {
		cfg.DB.User = dbUser
	}
	if cfg.Downloader.VideoPath == "" {
		cfg.Downloader.VideoPath = os.TempDir()
	}

	return cfg
}
