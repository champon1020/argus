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

var Logger LogHandler

func init() {
	Logger.NewLogger()
}

func (l *LogHandler) NewLogger() {
	logfileDir := os.Getenv("ARGUS_LOG_PATH")
	logfile, err := os.OpenFile(logfileDir+"debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("Unable to open log file: %s\n", err)
		l.SetOutput(os.Stdout)
	} else {
		l.SetOutput(io.Writer(logfile))
	}
	l.SetFlags(log.Ldate | log.Ltime)
}

func (l LogHandler) StackTrace() []string {
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

func (l LogHandler) Log() {
}

func (l LogHandler) ErrorLog(errs []Error) {
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
