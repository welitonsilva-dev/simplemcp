package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type Level int

const (
	DEBUG Level = iota
	INFO
	ERROR
)

var (
	currentLevel = DEBUG

	debugLog *log.Logger
	infoLog  *log.Logger
	errorLog *log.Logger
)

func Init(logDir string) error {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório de logs: %w", err)
	}

	logFile, err := os.OpenFile(
		filepath.Join(logDir, "humancli-server.log"),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return fmt.Errorf("erro ao abrir arquivo de log: %w", err)
	}

	multi := io.MultiWriter(os.Stdout, logFile)

	debugLog = log.New(multi, "DEBUG ", log.Ldate|log.Ltime)
	infoLog = log.New(multi, "INFO  ", log.Ldate|log.Ltime)
	errorLog = log.New(multi, "ERROR ", log.Ldate|log.Ltime)

	return nil
}

func caller() string {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return ""
	}
	parts := strings.Split(file, "/")
	if len(parts) > 2 {
		parts = parts[len(parts)-2:]
	}
	return fmt.Sprintf("[%s:%d]", strings.Join(parts, "/"), line)
}

func Debug(format string, args ...interface{}) {
	if currentLevel <= DEBUG && debugLog != nil {
		debugLog.Printf("%s %s", caller(), fmt.Sprintf(format, args...))
	}
}

func Info(format string, args ...interface{}) {
	if currentLevel <= INFO && infoLog != nil {
		infoLog.Printf("%s %s", caller(), fmt.Sprintf(format, args...))
	}
}

func Error(format string, args ...interface{}) {
	if currentLevel <= ERROR && errorLog != nil {
		errorLog.Printf("%s %s", caller(), fmt.Sprintf(format, args...))
	}
}

func SetLevel(l Level) {
	currentLevel = l
}
