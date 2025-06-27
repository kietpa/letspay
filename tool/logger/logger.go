package logger

import (
	"context"
	"io"
	"letspay/common/constants"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Config struct {
	FilePath   string `yaml:"file_path"` // "/var/log/app.log"
	MaxSizeMB  int    `yaml:"max_size_mb"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAgeDays int    `yaml:"max_age_days"`
	Compress   bool   `yaml:"compress"`

	LokiURL    string            `yaml:"loki_url"` // "http://loki:3100/loki/api/v1/push"
	LokiLabels map[string]string `yaml:"loki_labels"`
}

func New(cfg Config) zerolog.Logger {
	if err := os.MkdirAll(filepath.Dir(cfg.FilePath), 0755); err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}

	file, err := os.OpenFile(cfg.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	// lokiWriter, err := NewLokiWriter(cfg.LokiURL, cfg.LokiLabels)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixNano

	consoleWriter := zerolog.ConsoleWriter{
		Out: os.Stdout,
	}

	lumber := lumberjack.Logger{
		Filename:   cfg.FilePath,
		MaxSize:    cfg.MaxSizeMB,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAgeDays,
		Compress:   cfg.Compress,
	}

	multiWriter := io.MultiWriter(consoleWriter, &lumber, file)

	zerologger := zerolog.New(multiWriter).With().Timestamp().Logger()

	return zerologger
}

func Info(ctx context.Context, log zerolog.Logger, msg string) {
	pc, file, line, _ := runtime.Caller(2)

	funcName := getFuncName(runtime.FuncForPC(pc).Name())
	fileName := filepath.Base(file)
	processId, _ := ctx.Value(constants.PROCESS_ID).(string)

	log.Info().
		Str("file", fileName+":"+strconv.Itoa(line)).
		Str("func", funcName).
		Str(constants.PROCESS_ID, processId).
		Msg(msg)
}

func Error(ctx context.Context, log zerolog.Logger, msg string) {

	pc, file, line, _ := runtime.Caller(2)

	funcName := getFuncName(runtime.FuncForPC(pc).Name())
	fileName := filepath.Base(file)
	processId, _ := ctx.Value(constants.PROCESS_ID).(string)

	log.Error().
		Str("file", fileName+":"+strconv.Itoa(line)).
		Str("func", funcName).
		Str(constants.PROCESS_ID, processId).
		Msg(msg)
}

func getFuncName(s string) string {
	index := strings.LastIndex(s, ".")
	return s[index+1:]
}
