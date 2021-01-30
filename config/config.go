package config

import (
	"log"
	"os"
)

// Config represents the all the environmental variables that should be present
// on start up
type Config struct {
	Port           string
	RedirectURL    string
	Unsplash       UnsplashConfig
	Onedrive       OnedriveConfig
	OpenWeatherMap OpenWeatherMapConfig
	Dropbox        DropboxConfig
	GoogleDrive    GoogleDriveConfig
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

// OpenWeatherMapConfig represents OpenWeatherMap's configuration variables
type OpenWeatherMapConfig struct {
	AppID string
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
		OpenWeatherMap: OpenWeatherMapConfig{
			AppID: getEnv("OPENWEATHER_APPID", ""),
		},
		Dropbox: DropboxConfig{
			Key: getEnv("DROPBOX_KEY", ""),
		},
		GoogleDrive: GoogleDriveConfig{
			Key:    getEnv("GOOGLE_DRIVE_KEY", ""),
			Secret: getEnv("GOOGLE_DRIVE_SECRET", ""),
		},
	}

	return Conf
}

// getEnv reads an environment variable and returns it or returns a default
// value if the variable is optional. Otherwise, if a required variable is not
// set, the program will crash
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	if defaultVal == "" {
		log.Fatalf("%s has not been set in your ENV", key)
	}

	return defaultVal
}
