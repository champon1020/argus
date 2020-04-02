package argus

import (
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

// initialize function of top level
func init() {
	EnvVars = NewEnv()
	StdLogger = *NewStdLogger()
	Logger = *NewLogger()
	GlobalConfig = *NewConfig()
}

func NewLogger() *LogHandler {
	l := new(LogHandler)
	l.SetFlags(log.Ldate | log.Ltime)
	l.SetOutput(os.Stdout)
	return l
}

func NewStdLogger() *LogHandler {
	l := new(LogHandler)
	l.SetFlags(log.Ldate | log.Ltime)
	l.SetOutput(os.Stdout)
	return l
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
