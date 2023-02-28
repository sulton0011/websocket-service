package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

const (
	// DebugMode indicates service mode is debug.
	DebugMode = "debug"
	// TestMode indicates service mode is test.
	TestMode = "test"
	// ReleaseMode indicates service mode is release.
	ReleaseMode = "release"

	// Time allowed to write a message to the peer.
	WriteWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	PongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	PingPeriod = (PongWait * 9) / 10

	// Maximum message size allowed from peer.
	MaxMessageSize = 512

	// ReadBufferSize and WriteBufferSize specify I/O buffer sizes in bytes. If a buffer
	// size is zero, then buffers allocated by the HTTP server are used. The
	// I/O buffer sizes do not limit the size of the messages that can be sent
	// or received.
	WriteBufferSize = 1024
	ReadBufferSize  = 1024

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

	SecretKey string

	// InstrumetService
	InstrumentServiceHost string
	InstrumentGRPCPort    string
}

// Load ...
func Load() Config {

	envFileName := getOrReturnDefaultValue("ENV_FILE_PATH", "./.env").(string)
	if err := godotenv.Load(envFileName); err != nil {
		fmt.Println("No .env file found")
	}

	config := Config{}

	config.ServiceName = getOrReturnDefaultValue("SERVICE_NAME", "websocket_service").(string)
	config.Environment = getOrReturnDefaultValue("ENVIRONMENT", DebugMode).(string)
	config.Version = getOrReturnDefaultValue("VERSION", "1.0").(string)

	config.HTTPPort = getOrReturnDefaultValue("HTTP_PORT", ":8081").(string)
	config.HTTPScheme = getOrReturnDefaultValue("HTTP_SCHEME", "http").(string)

	config.SecretKey = getOrReturnDefaultValue("SECRET_KEY", "3Xn0zquCr4oO").(string)

	// Instrument Service
	config.InstrumentServiceHost = getOrReturnDefaultValue("INSTRUMENT_SERVICE_HOST", "0.0.0.0").(string)
	config.InstrumentGRPCPort = getOrReturnDefaultValue("INSTRUMENT_GRPC_PORT", ":9103").(string)

	return config
}

func getOrReturnDefaultValue(key string, defaultValue interface{}) interface{} {
	val, exists := os.LookupEnv(key)

	if exists {
		return val
	}

	return defaultValue
}
