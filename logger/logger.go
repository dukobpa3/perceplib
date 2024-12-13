package logger

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Кешуємо часто використовувані поля
var commonFields = map[string]struct{}{
	"error":   {},
	"service": {},
	"time":    {},
	"level":   {},
}

func NewLogger(level LogLevel, options ...Option) *Logger {
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:  "message",
		LevelKey:    "level",
		TimeKey:     "time",
		EncodeLevel: zapcore.CapitalColorLevelEncoder,
		EncodeTime:  zapcore.ISO8601TimeEncoder,
	}

	config := zap.Config{
		Level:         zap.NewAtomicLevelAt(zapcore.Level(level)),
		Development:   true,
		Encoding:      "console",
		OutputPaths:   []string{"stdout"},
		EncoderConfig: encoderConfig,
	}

	zapLogger, _ := config.Build(convertOptions(options)...)
	return &Logger{
		zapLogger: zapLogger,
		services:  make(map[string]ServiceConfig),
		mu:        &sync.RWMutex{},
	}
}

func (l *Logger) log(level LogLevel, msg string, fields ...Field) {
	if len(l.zapLogger.Name()) > 0 && !l.isServiceEnabled(l.zapLogger.Name()) {
		return
	}

	zapFields := convertFields(fields)
	switch level {
	case DebugLevel:
		l.zapLogger.Debug(msg, zapFields...)
	case InfoLevel:
		l.zapLogger.Info(msg, zapFields...)
	case WarnLevel:
		l.zapLogger.Warn(msg, zapFields...)
	case ErrorLevel:
		l.zapLogger.Error(msg, zapFields...)
	case DPanicLevel:
		l.zapLogger.DPanic(msg, zapFields...)
	case PanicLevel:
		l.zapLogger.Panic(msg, zapFields...)
	case FatalLevel:
		l.zapLogger.Fatal(msg, zapFields...)
	}
}

func (l *Logger) Debug(msg string, fields ...Field) {
	l.log(DebugLevel, msg, fields...)
}

func (l *Logger) Info(msg string, fields ...Field) {
	l.log(InfoLevel, msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...Field) {
	l.log(WarnLevel, msg, fields...)
}

func (l *Logger) Error(msg string, fields ...Field) {
	l.log(ErrorLevel, msg, fields...)
}

func (l *Logger) DPanic(msg string, fields ...Field) {
	l.log(DPanicLevel, msg, fields...)
}

func (l *Logger) Panic(msg string, fields ...Field) {
	l.log(PanicLevel, msg, fields...)
}

func (l *Logger) Fatal(msg string, fields ...Field) {
	l.log(FatalLevel, msg, fields...)
}

func (l *Logger) Named(name string) *Logger {
	color := getColorForService(name)
	if _, exists := l.services[name]; !exists {
		l.RegisterService(name, color)
	}

	namedLogger := newColoredLogger(name, color)
	return &Logger{
		zapLogger: namedLogger,
		services:  l.services,
		mu:        l.mu,
	}
}

func (l *Logger) With(fields ...Field) *Logger {
	zapFields := convertFields(fields)
	return &Logger{
		zapLogger: l.zapLogger.With(zapFields...),
		services:  l.services,
		mu:        l.mu,
	}
}

func (l *Logger) WithOptions(options ...Option) *Logger {
	zapOptions := convertOptions(options)
	return &Logger{
		zapLogger: l.zapLogger.WithOptions(zapOptions...),
		services:  l.services,
		mu:        l.mu,
	}
}

func convertOptions(options []Option) []zap.Option {
	zapOptions := make([]zap.Option, len(options))
	for i, opt := range options {
		zapOptions[i] = opt.zapOption
	}
	return zapOptions
}

func convertFields(fields []Field) []zap.Field {
	if len(fields) == 0 {
		return nil
	}

	zapFields := make([]zap.Field, len(fields))
	for i, field := range fields {
		if _, ok := commonFields[field.Key]; ok {
			zapFields[i] = zap.Any(field.Key, field.Value)
			continue
		}
		zapFields[i] = zap.Any(field.Key, field.Value)
	}
	return zapFields
}

func convertZapFields(zapFields []zap.Field) []Field {
	fields := make([]Field, len(zapFields))
	for i, zapField := range zapFields {
		fields[i] = Field{
			Key:   zapField.Key,
			Value: zapField.Interface,
		}
	}
	return fields
}
