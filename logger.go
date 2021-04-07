package argus

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
)

// Logger object.
type Logger struct {
	TimeFormat string
}

// NewLogger creates Logger instance.
func NewLogger() *Logger {
	return &Logger{TimeFormat: time.RFC3339}
}

func (l *Logger) Error(c echo.Context, status int, err error) {
	// TIME | METHOD | REMOTE_IP | URI | STATUS | MESSAGE
	fmt.Printf("[ARGUS] %v | %v | %s | %s | %d | %s\n",
		time.Now().Format(l.TimeFormat),
		c.Request().Method,
		c.RealIP(),
		c.Path(),
		status,
		err.Error())
}
