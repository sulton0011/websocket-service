package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

const (
	// DebugMode indicates service mode is debug.
	DebugMode = "debug"
	// TestMode indicates service mode is test.
	TestMode = "test"
	// ReleaseMode indicates service mode is release.
	ReleaseMode = "release"

	// Errors mode
	ErrorModel = "!!!Error"
	ErrorStyle = "-->"
)

type Config struct {
	ServiceName string
	Environment string // debug, test, release
	Version     string

	// HTTP
	HTTPPort   string
	HTTPScheme string

	// InstrumetService
	InstrumentServiceHost string
	InstrumentGRPCPort    string
}

// Load ...
func Load() Config {

	envFileName := getOrReturnDefaultValue("ENV_FILE_PATH", "../.env").(string)
	if err := godotenv.Load(envFileName); err != nil {
		fmt.Println("No .env file found")
	}

	config := Config{}

	config.ServiceName = getOrReturnDefaultValue("SERVICE_NAME", "hyssa_go_api_gateway").(string)
	config.Environment = getOrReturnDefaultValue("ENVIRONMENT", DebugMode).(string)
	config.Version = getOrReturnDefaultValue("VERSION", "1.0").(string)

	config.HTTPPort = getOrReturnDefaultValue("HTTP_PORT", ":8081").(string)
	config.HTTPScheme = getOrReturnDefaultValue("HTTP_SCHEME", "http").(string)

	// Instrument Service
	config.InstrumentServiceHost = getOrReturnDefaultValue("INSTRUMENT_SERVICE_HOST", "0.0.0.0").(string)
	config.InstrumentGRPCPort = getOrReturnDefaultValue("INSTRUMENT_GRPC_PORT", ":9103").(string)

	return Config{}
}

func getOrReturnDefaultValue(key string, defaultValue interface{}) interface{} {
	val, exists := os.LookupEnv(key)

	if exists {
		return val
	}

	return defaultValue
}
