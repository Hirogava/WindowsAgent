package services

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

var (
	frontendLogger     *log.Logger
	frontendLoggerOnce sync.Once
)

func getFrontendLogger() *log.Logger {
	frontendLoggerOnce.Do(func() {
		cs := NewConfigService()
		rootDir, err := cs.repoRootDir()
		if err != nil {
			frontendLogger = log.New(os.Stdout, "[frontend] ", log.LstdFlags|log.Lshortfile)
			frontendLogger.Printf("logger init warning: %v", err)
			return
		}

		logsDir := filepath.Join(rootDir, "logs")
		if mkErr := os.MkdirAll(logsDir, 0755); mkErr != nil {
			frontendLogger = log.New(os.Stdout, "[frontend] ", log.LstdFlags|log.Lshortfile)
			frontendLogger.Printf("logger init warning: %v", mkErr)
			return
		}

		logPath := filepath.Join(logsDir, "frontend.log")
		file, openErr := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if openErr != nil {
			frontendLogger = log.New(os.Stdout, "[frontend] ", log.LstdFlags|log.Lshortfile)
			frontendLogger.Printf("logger init warning: %v", openErr)
			return
		}

		writer := io.MultiWriter(os.Stdout, file)
		frontendLogger = log.New(writer, "[frontend] ", log.LstdFlags|log.Lshortfile)
		frontendLogger.Printf("logging initialized: %s", fmt.Sprintf("%s", logPath))
	})

	return frontendLogger
}
