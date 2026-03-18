package server

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"
)

const logDir = "log"

var logCounter uint64

// LogResponseToDir writes JSON response data to a timestamped file in the given directory.
func LogResponseToDir(dir string, data []byte) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Printf("Error creating log directory: %v", err)
		return
	}
	seq := atomic.AddUint64(&logCounter, 1)
	filename := filepath.Join(dir, fmt.Sprintf("%d_%d.json", time.Now().UnixNano(), seq))
	if err := os.WriteFile(filename, data, 0644); err != nil {
		log.Printf("Error writing log file: %v", err)
	}
}

// LogResponse writes JSON response data to the default log directory.
func LogResponse(data []byte) {
	LogResponseToDir(logDir, data)
}
