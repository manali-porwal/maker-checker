package config

import (
	"github.com/spf13/viper"
)

var cfg *AppConfig

type AppConfig struct {
	App struct {
		Name  string
		Debug bool
		Env   string
		Port  uint16
	}

	Gin struct {
		Mode string
	}

	Database struct {
		ReadDSN  string
		WriteDSN string
	}

	Checker struct {
		NumRequiredApprovals int
	}

	JWT struct {
		Secret      string
		ExpiryHours int
	}
}

func Config() *AppConfig {
	if cfg == nil {
		loadConfig()
	}

	return cfg
}

func loadConfig() {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	// Ignore config file not found, perhaps we will use environment variables.
	_ = viper.ReadInConfig()

	cfg = &AppConfig{}

	// App.
	cfg.App.Name = viper.GetString("APP_NAME")
	cfg.App.Debug = viper.GetBool("APP_DEBUG")
	cfg.App.Env = viper.GetString("APP_ENV")
	cfg.App.Port = viper.GetUint16("APP_PORT")

	// Gin.
	cfg.Gin.Mode = viper.GetString("GIN_MODE")

	// Database.
	cfg.Database.ReadDSN = viper.GetString("DB_READ_DSN")
	cfg.Database.WriteDSN = viper.GetString("DB_WRITE_DSN")

	// Checker.
	cfg.Checker.NumRequiredApprovals = viper.GetInt("NUM_REQUIRED_APPROVALS")

	// JWT.
	cfg.JWT.Secret = viper.GetString("JWT_SECRET")
	cfg.JWT.ExpiryHours = viper.GetInt("JWT_EXPIRY")
}
