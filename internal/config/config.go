package config

import (
	"flag" // Used to read command-line flags like: --config=path
	"log"  // For logging errors and fatal messages
	"os"   // For environment variables & file system checks

	"github.com/ilyakaznacheev/cleanenv" // Library to read YAML/env config easily
)

// HTTPServer holds the HTTP server address from config
type HTTPServer struct {
	Addr string `yaml:"address" env-required:"true"` // Server address (e.g., :8082)
}

// Config struct maps YAML + ENV variables for the entire app
type Config struct {
	Env         string               `yaml:"env" env:"ENV" env-required:"true"` // App environment (dev/prod)
	StoragePath string               `yaml:"storage_path" env-required:"true"`  // Path to database/storage file
	HTTPServer  `yaml:"http_server"` // Embedded struct for server config
}

func MustLoad() *Config {
	var configPath string

	// First: Check if CONFIG_PATH environment variable is set
	configPath = os.Getenv("CONFIG_PATH")

	// If CONFIG_PATH is not set, fallback to command-line flag: --config=path
	if configPath == "" {
		flags := flag.String("config", "", "path to the configuration file")
		flag.Parse()

		configPath = *flags

		if configPath == "" {
			log.Fatal("Config path is not set")
		}
	}

	// Check if config file actually exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist %s", configPath) // File missing → exit
	}

	var cfg Config // Object where YAML data will be stored

	// Read YAML config using cleanenv into cfg struct
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("can't read config file: %s", err.Error()) // Config parsing failed → exit
	}

	return &cfg // Return loaded configuration
}
