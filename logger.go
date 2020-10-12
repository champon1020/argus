package argus

import (
	"log"
	"os"
)

// LogHandler handles application log processes.
type LogHandler struct {
	log.Logger
}

// NewLogger creates logger instance.
func NewLogger() *LogHandler {
	l := new(LogHandler)
	l.SetFlags(log.Ldate | log.Ltime)
	l.SetOutput(os.Stdout)
	return l
}
