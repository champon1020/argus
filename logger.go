package argus

import (
	"io"
	"log"
	"os"
	"runtime"
	"strconv"
)

type LogHandler struct {
	log.Logger
}

var (
	Logger    LogHandler
	StdLogger LogHandler
)

func (l *LogHandler) New() {
	logfileDir := EnvVars.Get("log")
	logfile, err := os.OpenFile(logfileDir+"debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	l.SetFlags(log.Ldate | log.Ltime)
	if err != nil {
		log.Printf("Unable to open log file: %s\n", err)
		l.SetOutput(os.Stdout)
		return
	}
	l.SetOutput(io.Writer(logfile))
}

func (l *LogHandler) NewStd() {
	l.SetFlags(log.Ldate | log.Ltime)
	l.SetOutput(os.Stdout)
}

func (l *LogHandler) StackTrace() []string {
	c := 0
	var stackTrace []string
	for {
		if pc, file, line, ok := runtime.Caller(c); ok {
			funcName := runtime.FuncForPC(pc).Name()
			st := "func=" + funcName + ", line=" + strconv.Itoa(line) + ", file=" + file + "\n"
			stackTrace = append(stackTrace, st)
			c++
			continue
		}
		break
	}
	return stackTrace
}

func (l *LogHandler) Log() {
}

func (l *LogHandler) ErrorLog(errs []Error) {
	if len(errs) == 0 {
		return
	}
	var rows []string
	for _, v := range errs {
		j, ok := v.Marshal()
		if ok != nil {
			Logger.Fatalf("Unable to parse error to json\n")
		}
		rows = append(rows, string(j))
	}

	jsonMap := map[string]interface{}{
		"Errors":     rows,
		"StackTrace": l.StackTrace(),
	}
	l.Println(jsonMap)
}
