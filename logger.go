package argus

import (
	"io"
	"log"
	"os"
	"runtime"
)

type Logger struct {
	log.Logger
}

var logger Logger

func init() {
	logger.NewLogger("[argus]")
}

func (l *Logger) NewLogger(prefix string) {
	logfileDir := os.Getenv("ARGUS_LOG_PATH")
	logfile, err := os.OpenFile(logfileDir+"debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("Unable to open log file: %s\n", err)
		l.SetOutput(os.Stdout)
	} else {
		l.SetOutput(io.Writer(logfile))
	}
	l.SetPrefix(prefix)
	l.SetFlags(log.Ldate | log.Ltime)
}

func (l Logger) StackTrace() {
	c := 0
	for {
		if pc, file, line, ok := runtime.Caller(c); ok {
			funcName := runtime.FuncForPC(pc).Name()
			l.Printf("func=%v, line=%d, file=%s", funcName, line, file)
			c++
			continue
		}
		break
	}
}

func (l Logger) ErrorPrintf(err error) {
	l.StackTrace()
	l.Printf("%v\n", err)
}

func (l Logger) ErrorMsgPrintf(msg string, err error) {
	l.StackTrace()
	l.Printf("%s: %v\n", msg, err)
}

func (l Logger) ErrorFatalf(err error) {
	l.StackTrace()
	l.Fatalf("%v\n", err)
}

func (l Logger) ErrorPanic(err error) {
	l.StackTrace()
	panic(err)
}
