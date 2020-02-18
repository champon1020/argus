package argus

import (
	"io"
	"log"
	"os"
)

type Logger struct {
	log.Logger
}

var logger Logger

func init() {
	logger.NewLogger("[argus]")
}

func (l *Logger) NewLogger(prefix string) {
	logfile, err := os.OpenFile("logs/debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("Unable to open log file: %s\n", err)
		l.SetOutput(os.Stdout)
	} else {
		l.SetOutput(io.Writer(logfile))
	}
	l.SetPrefix(prefix)
	l.SetFlags(log.Ldate | log.Ltime)
}

func (l Logger) ErrorPrintf(err error) {
	l.Printf("%s\n", err)
}

func (l Logger) ErrorMsgPrintf(msg string, err error) {
	l.Printf(msg+": %s\n", err)
}

func (l Logger) ErrorFatalf(err error) {
	l.Fatalf("%s\n", err)
}

func (l Logger) ErrorPanic(err error) {
	panic(err)
}
