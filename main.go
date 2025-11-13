package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"

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

var (
	config    Config
	configMu  sync.RWMutex
	buildTime string // set at compile time via -ldflags
)

func init() {
	if buildTime == "" {
		buildTime = "unknown"
	}
	log.Printf("Go Version: %s", runtime.Version())
	log.Printf("Build Time: %s", buildTime)
}

func main() {
	// Load configuration
	if err := loadConfig(); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Set up signal handling for config reload
	setupSignalHandling()

	// Register the handler for the root path
	http.HandleFunc("/", handleRoot)

	// Start the server
	configMu.RLock()
	addr := fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)
	message := config.Message
	configMu.RUnlock()

	log.Printf("Starting HTTP server on http://%s\n", addr)
	log.Printf("Response message: %s\n", message)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// setupSignalHandling sets up signal handlers for graceful reload
func setupSignalHandling() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP)

	go func() {
		for {
			sig := <-sigs
			if sig == syscall.SIGHUP {
				log.Println("Received SIGHUP, reloading configuration...")
				if err := loadConfig(); err != nil {
					log.Printf("Failed to reload configuration: %v", err)
				} else {
					configMu.RLock()
					log.Printf("Configuration reloaded successfully. New message: %s", config.Message)
					configMu.RUnlock()
				}
			}
		}
	}()
}

// loadConfig loads configuration from file
func loadConfig() error {
	// Try multiple config file locations
	configPaths := []string{
		"/etc/dpkg-build-pg/config.production.yaml", // Production config (highest priority)
		"/etc/dpkg-build-pg/config.yaml",            // System-wide config
		"./config.yaml",                             // Local config for development
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

	// Parse YAML into temporary config
	var newConfig Config
	if err := yaml.Unmarshal(configData, &newConfig); err != nil {
		return fmt.Errorf("failed to parse config file: %v", err)
	}

	// Set defaults if not specified
	if newConfig.Server.Port == 0 {
		newConfig.Server.Port = 8080
	}
	if newConfig.Server.Host == "" {
		newConfig.Server.Host = "0.0.0.0"
	}
	if newConfig.Message == "" {
		newConfig.Message = "Hello from dpkg-build-pg server!"
	}

	// Update global config with lock
	configMu.Lock()
	config = newConfig
	configMu.Unlock()

	return nil
}

// handleRoot is the handler function that responds with a string from config
func handleRoot(w http.ResponseWriter, r *http.Request) {
	configMu.RLock()
	message := config.Message
	configMu.RUnlock()
	fmt.Fprintf(w, message)
}
