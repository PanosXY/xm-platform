package config

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"strconv"
)

const (
	appName      = "xm-service"
	profilesPath = "profiles"
)

var (
	//go:embed profiles/*
	configs embed.FS
)

// Configuration is the service's configuration
type Configuration struct {
	App        *AppConfig
	HttpServer *HttpServerConfig
	Database   *DatabaseConfig
}

// AppConfig is the application's configuration
type AppConfig struct {
	Env     string
	Version string
	Name    string
}

// HttpServerConfig is the server's configuration
type HttpServerConfig struct {
	Host         string
	Port         int
	JWTSecretKey []byte
}

// DatabaseConfig is the database's configuration
type DatabaseConfig struct {
	Username string
	Password string
	Name     string
	SSL      string
	Host     string
	Port     int
	Settings *DatabaseSettings
}

type ProfileConfig struct {
	Database DatabaseSettings `json:"database"`
}

type DatabaseSettings struct {
	ConnMaxLifetimeSecs int `json:"max_open_connections"`
	ConnMaxIdleTimeSecs int `json:"max_idle_connections"`
	MaxOpenConns        int `json:"connection_max_lifetime_secs"`
	MaxIdleConns        int `json:"connection_max_idle_time_secs"`
	ConnectTimeout      int `json:"connect_timeout_secs"`
}

// LoadConfiguration returns the service's configuration
func LoadConfiguration() (*Configuration, error) {
	httpServerConfig, err := getHttpServerConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get http server configuration: %v", err)
	}

	profile, err := loadConfigurationProfile(os.Getenv("APP_ENV"))
	if err != nil {
		return nil, fmt.Errorf("failed to get profile configuration")
	}

	databaseConfig, err := getDatabaseConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get database configuration: %v", err)
	}
	databaseConfig.Settings = &profile.Database

	return &Configuration{
		App:        getAppConfig(),
		HttpServer: httpServerConfig,
		Database:   databaseConfig,
	}, nil
}

func getAppConfig() *AppConfig {
	return &AppConfig{
		Env:     os.Getenv("APP_ENV"),
		Version: os.Getenv("APP_VERSION"),
		Name:    appName,
	}
}

func getHttpServerConfig() (*HttpServerConfig, error) {
	port, err := strconv.Atoi(os.Getenv("APP_PORT"))
	if err != nil {
		return nil, err
	}

	return &HttpServerConfig{
		Host:         os.Getenv("APP_HOST"),
		Port:         port,
		JWTSecretKey: []byte(os.Getenv("JWT_SECRET_KEY")),
	}, nil
}

func getDatabaseConfig() (*DatabaseConfig, error) {
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return nil, err
	}

	return &DatabaseConfig{
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		Name:     os.Getenv("DB_NAME"),
		SSL:      os.Getenv("DB_SSL"),
		Host:     os.Getenv("DB_HOST"),
		Port:     port,
	}, nil
}

func loadConfigurationProfile(environment string) (*ProfileConfig, error) {
	files, err := fs.ReadDir(configs, profilesPath)
	if err != nil {
		return nil, err
	}

	var config *ProfileConfig

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := fmt.Sprintf("%s.json", environment)
		if file.Name() != fmt.Sprintf("%s.json", environment) {
			continue
		}

		configFile, err := configs.ReadFile(fmt.Sprintf("%s/%s", profilesPath, filePath))
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(configFile, &config); err != nil {
			return nil, err
		}

		return config, nil
	}

	return nil, fmt.Errorf("failed to find configuration environment: %s", environment)
}
