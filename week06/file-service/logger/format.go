package logger

import (
	"time"

	"github.com/gin-gonic/gin"
)

type LogEntry struct {
	Timestamp   string      `json:"timestamp"`
	HTTPMethod  string      `json:"http_method"`
	Path        string      `json:"path"`
	StatusCode  int         `json:"status_code"`
	RequestBody interface{} `json:"request_body,omitempty"`
	Error       string      `json:"error,omitempty"`
}

func NewLogEntry(c *gin.Context) *LogEntry {
	return &LogEntry{
		Timestamp:  time.Now().Format(time.RFC3339Nano),
		HTTPMethod: c.Request.Method,
		Path:       c.Request.URL.Path,
	}
}
