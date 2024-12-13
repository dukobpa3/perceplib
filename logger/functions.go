package logger

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	serviceColors = make(map[string]string)
	colors        = []string{
		ColorRed, ColorGreen, ColorYellow, ColorBlue,
		ColorMagenta, ColorCyan, ColorWhite, ColorBrightBlack,
		ColorBrightRed, ColorBrightGreen, ColorBrightYellow,
		ColorBrightBlue, ColorBrightMagenta, ColorBrightCyan,
		ColorBrightWhite,
	}
	globalMu sync.RWMutex
	r        = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// Helper functions for creating Fields
func String(key, val string) Field {
	return Field{Key: key, Value: val}
}

func Any(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}

func Error(err error) Field {
	return Field{Key: "error", Value: err}
}

func NewDevelopment() (*Logger, error) {
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	return &Logger{
		zapLogger: zapLogger,
		services:  make(map[string]ServiceConfig),
		mu:        &sync.RWMutex{},
	}, nil
}

func Level(level LogLevel) LogLevel {
	return LogLevel(zapcore.Level(level))
}

func IncreaseLevel(level LogLevel) Option {
	return Option{zapOption: zap.IncreaseLevel(zapcore.Level(level))}
}

func getColorForService(service string) string {
	globalMu.Lock()
	defer globalMu.Unlock()

	if color, exists := serviceColors[service]; exists {
		return color
	}

	color := colors[r.Intn(len(colors))]
	serviceColors[service] = color
	return color
}

// Logger creation
func newColoredLogger(service string, color string) *zap.Logger {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		zapcore.DebugLevel,
	)

	logger := zap.New(core).With(zap.String("service", fmt.Sprintf("%s%s%s", color, service, ColorReset)))

	return logger
}
