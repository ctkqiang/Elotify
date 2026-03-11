package config

import (
	"os"
	"path/filepath"
	"pushnotification_services/internal/utilities"

	"github.com/joho/godotenv"
)

type LogMode string

var (
	LOG_DEBUG   = LogMode("DEBUG")
	LOG_VERBOSE = LogMode("VERBOSE")
	LOG_INFO    = LogMode("INFO")
)

type MongoDBConfig struct {
	DatabaseUser     string `json:"database_user"`
	DatabasePassword string `json:"database_password"`
	DatabaseName     string `json:"database_name"`
	DatabaseHost     string `json:"database_host"`
	DatabasePort     string `json:"database_port"`
}

type OneSignalConfig struct {
	AppID  string
	APIKey string
}

type JWEConfig struct {
	KeyAES256 string `json:"key_aes_256"`
	Salt      string `json:"salt"`
	Cipher    string `json:"cipher"`
	IsCaesar  bool   `json:"is_caesar"`
	Encrypt   bool   `json:"encrypt"`
}

var (
	COLLECTION_NOTIFICATIONS = "notifications"
	COLLECTION_ANNOUNCEMENTS = "announcements"
	COLLECTION_ADVERTISEMENTS = "advertisements"
)

var (
	MongoDBCreds   MongoDBConfig
	OneSignalCreds OneSignalConfig
	JWECreds       JWEConfig
)

func init() {
	selectedEnviornmentFlavours, err := GetDevelopmentFlavours()
	if err != nil {
		utilities.Log(utilities.ERROR, "Environment Error %s", err.Error())
	}

	workspaceDirectory, err := os.Getwd()
	if err == nil {
		envPath := filepath.Join(workspaceDirectory, "internal", "config", selectedEnviornmentFlavours)

		if _, err := os.Stat(envPath); err == nil {
			if err := godotenv.Load(envPath); err != nil {
				utilities.Log(utilities.ERROR, "[CONFIG] WARNING!!! Failed to load [.env] File %s: %v\n", envPath, err)
			}
		} else {
			godotenv.Load()
		}
	}

	// Docker环境检测和主机配置
	host := os.Getenv("MONGODB_HOST")
	if host == "localhost" || host == "" {
		if os.Getenv("DOCKER_ENV") == "true" {
			host = "mongodb"
		}
	}

	MongoDBCreds = MongoDBConfig{
		DatabaseUser:     os.Getenv("MONGODB_USER"),
		DatabasePassword: os.Getenv("MONGODB_PASSWORD"),
		DatabaseName:     os.Getenv("MONGODB_DATABASE"),
		DatabaseHost:     host,
		DatabasePort:     os.Getenv("MONGODB_PORT"),
	}

	OneSignalCreds = OneSignalConfig{
		AppID:  os.Getenv("ONESIGNAL_APP_ID"),
		APIKey: os.Getenv("ONESIGNAL_REST_API_KEY"),
	}

	JWECreds = JWEConfig{
		KeyAES256: os.Getenv("JWE_KEY_AES_256"),
		Salt:      "no-salt",
		Cipher:    "AES-256-CBC",
		IsCaesar:  false,
		Encrypt:   os.Getenv("ENCRYPT") == "true",
	}
}