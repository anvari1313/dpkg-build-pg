package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Server struct {
		Port int    `yaml:"port"`
		Host string `yaml:"host"`
	} `yaml:"server"`
	Message string `yaml:"message"`
}

var config Config

// buildTime is set at compile time via -ldflags
var buildTime string

func init() {
	if buildTime == "" {
		buildTime = "unknown"
	}
	log.Printf("Build Time: %s", buildTime)
}

func main() {
	// Load configuration
	if err := loadConfig(); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Register the handler for the root path
	http.HandleFunc("/", handleRoot)

	// Start the server
	addr := fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)
	log.Printf("Starting HTTP server on http://%s\n", addr)
	log.Printf("Response message: %s\n", config.Message)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// loadConfig loads configuration from file
func loadConfig() error {
	// Try multiple config file locations
	configPaths := []string{
		"/etc/dpkg-build-pg/config.yaml", // System-wide config
		"./config.yaml",                  // Local config for development
	}

	var configData []byte
	var err error
	var usedPath string

	for _, path := range configPaths {
		configData, err = os.ReadFile(path)
		if err == nil {
			usedPath = path
			break
		}
	}

	if err != nil {
		return fmt.Errorf("could not find config file in any of the expected locations: %v", configPaths)
	}

	log.Printf("Loading config from: %s", usedPath)

	// Parse YAML
	if err := yaml.Unmarshal(configData, &config); err != nil {
		return fmt.Errorf("failed to parse config file: %v", err)
	}

	// Set defaults if not specified
	if config.Server.Port == 0 {
		config.Server.Port = 8080
	}
	if config.Server.Host == "" {
		config.Server.Host = "0.0.0.0"
	}
	if config.Message == "" {
		config.Message = "Hello from dpkg-build-pg server!"
	}

	return nil
}

// handleRoot is the handler function that responds with a string from config
func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, config.Message)
}
