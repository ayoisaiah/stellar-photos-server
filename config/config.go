package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/ayoisaiah/stellar-photos-server/utils"
)

// Config represents the all the environmental variables that should be present
// on start up
type Config struct {
	Port        string
	RedirectURL string
	Unsplash    UnsplashConfig
	Onedrive    OnedriveConfig
	Dropbox     DropboxConfig
	GoogleDrive GoogleDriveConfig
	Redis       RedisConfig
	LogLevel    int
}

// RedisConfig represents redis configuration variables
type RedisConfig struct {
	Addr     string
	DB       int
	Password string
	Username string
}

// UnsplashConfig represents Unsplash's API configuration variables
type UnsplashConfig struct {
	AccessKey string
}

// OnedriveConfig represents Onedrive's API configuration variables
type OnedriveConfig struct {
	AppID  string
	Secret string
}

// DropboxConfig represents Dropbox's API configuration variables
type DropboxConfig struct {
	Key string
}

// GoogleDriveConfig represents Google Drive's API configuration variables
type GoogleDriveConfig struct {
	Key    string
	Secret string
}

// Conf represents the application configuration
var Conf *Config

// New returns a new Config struct
func New() *Config {
	redisDBStr := getEnv("REDIS_DB", "")
	redisDBInt, err := strconv.Atoi(redisDBStr)
	if err != nil || redisDBInt < 0 {
		utils.Logger().
			Fatalw("ENV: REDIS_DB must be a positive integer", "tag", "redis_db_env", "REDIS_DB", redisDBStr)
	}

	logLevel := getEnv("LOG_LEVEL", "0")
	logLevelInt, err := strconv.Atoi(logLevel)
	if err != nil {
		utils.Logger().
			Fatalw("ENV: LOG_LEVEL must be an integer", "tag", "log_level_env", "REDIS_DB", redisDBStr)
	}

	Conf = &Config{
		Port:        getEnv("PORT", "8080"),
		RedirectURL: getEnv("REDIRECT_URL", ""),
		Unsplash: UnsplashConfig{
			AccessKey: getEnv("UNSPLASH_ACCESS_KEY", ""),
		},
		Onedrive: OnedriveConfig{
			AppID:  getEnv("ONEDRIVE_APPID", ""),
			Secret: getEnv("ONEDRIVE_SECRET", ""),
		},
		Dropbox: DropboxConfig{
			Key: getEnv("DROPBOX_KEY", ""),
		},
		GoogleDrive: GoogleDriveConfig{
			Key:    getEnv("GOOGLE_DRIVE_KEY", ""),
			Secret: getEnv("GOOGLE_DRIVE_SECRET", ""),
		},
		Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
			Username: getEnv("REDIS_USERNAME", ""),
			DB:       redisDBInt,
			Password: getEnv("REDIS_PASSWORD", ""),
		},
		LogLevel: logLevelInt,
	}

	return Conf
}

// getEnv reads an environment variable and returns it or returns a default
// value if the variable is optional. Otherwise, if a required variable is not
// set, the program will crash
func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	if defaultVal == "" {
		utils.Logger().
			Fatalw(fmt.Sprintf("%s has not been set in your ENV", key), "tag", "env_not_set", "key")
	}

	return defaultVal
}
