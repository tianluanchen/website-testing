package pkg

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/mattn/go-isatty"
)

// log level
const (
	LDEBUG = iota + 1 // 1
	LINFO             // 2
	LWARN             // 3
	LERROR            // 4
	LFATAL            // 5
)

var logMutex sync.Mutex
var logOutput io.Writer
var logLevel int = LDEBUG
var logColor bool = true
var logOutputIsTerm bool

func init() {
	SetLogOutput(os.Stderr)
}

func doWithLock(f func()) {
	logMutex.Lock()
	defer logMutex.Unlock()
	f()
}

// if level >= logLevel, continue
func checkLogLevel(level int) bool {
	return level >= logLevel
}

func SetLogOutput(w io.Writer) {
	doWithLock(func() {
		logOutput = w
		if w, ok := logOutput.(*os.File); !ok || os.Getenv("TERM") == "dumb" ||
			(!isatty.IsTerminal(w.Fd()) && !isatty.IsCygwinTerminal(w.Fd())) {
			logOutputIsTerm = false
		} else {
			logOutputIsTerm = true
		}
	})
}

// if v is false, disable output color
func SetLogColor(v bool) {
	doWithLock(func() {
		logColor = v
	})
}

// use LDEBUG, LINFO, LWARN, LERROR, LFATAL
func SetLogLevel(level int) {
	doWithLock(func() {
		logLevel = level
	})
}

type Logger struct {
	prefix string
}

func NewLogger(prefix ...string) *Logger {
	l := &Logger{}
	if len(prefix) > 0 {
		l.prefix = prefix[0]
	}
	return l
}

// if format is not empty, similar to fmt.Printf
//
// if format is empty and ln is false, similar to fmt.Print
//
// if format is empty and ln is true, similar to fmt.Println
func (l *Logger) write(level int, format string, ln bool, v ...any) {
	if !checkLogLevel(level) {
		return
	}
	buf := &strings.Builder{}

	var (
		green        string
		white        string
		gray         string
		yellow       string
		red          string
		whiteFgRedBg string
		blue         string
		// magenta      string
		cyan  string
		reset string
	)

	if logColor && logOutputIsTerm {
		green = "\033[32m"
		white = "\033[97m"
		gray = "\033[90m"
		yellow = "\033[93m"
		red = "\033[91m"
		whiteFgRedBg = "\033[97;41m"
		blue = "\033[94m"
		// magenta = "\033[35m"
		cyan = "\033[36m"
		reset = "\033[0m"
	}
	buf.WriteString(green + time.Now().Format("2006-01-02T15:04:05.000Z0700") + reset + "  ")

	switch level {
	case LDEBUG:
		buf.WriteString(blue + "DEBUG" + reset)
	case LINFO:
		buf.WriteString(white + "INFO" + reset)
	case LWARN:
		buf.WriteString(yellow + "WARN" + reset)
	case LERROR:
		buf.WriteString(red + "ERROR" + reset)
	case LFATAL:
		buf.WriteString(whiteFgRedBg + "FATAL" + reset)
	default:
		buf.WriteString(gray + "UNKNOWN" + reset)
	}
	buf.WriteString("  ")
	_, file, line, _ := runtime.Caller(2)
	// pc, file, line, _ := runtime.Caller(2)
	// funcName := runtime.FuncForPC(pc).Name()
	index := strings.LastIndexByte(file, '/')
	if index > -1 {
		file = file[index+1:]
	}
	buf.WriteString(cyan + file + reset + ":" + cyan + strconv.Itoa(line) + reset + "  ")
	// buf.WriteString(magenta + l.prefix + reset)
	buf.WriteString(l.prefix)
	if format == "" {
		if ln {
			fmt.Fprintln(buf, v...)
		} else {
			fmt.Fprint(buf, v...)
			buf.WriteByte('\n')
		}
	} else {
		fmt.Fprintf(buf, format, v...)
		buf.WriteByte('\n')
	}
	doWithLock(func() {
		fmt.Fprint(logOutput, buf.String())
	})
}

func (l *Logger) Debug(v ...interface{}) {
	l.write(LDEBUG, "", false, v...)
}
func (l *Logger) Debugln(v ...interface{}) {
	l.write(LDEBUG, "", true, v...)
}
func (l *Logger) Debugf(format string, v ...interface{}) {
	l.write(LDEBUG, format, false, v...)
}
func (l *Logger) Info(v ...interface{}) {
	l.write(LINFO, "", false, v...)
}
func (l *Logger) Infoln(v ...interface{}) {
	l.write(LINFO, "", true, v...)
}
func (l *Logger) Infof(format string, v ...interface{}) {
	l.write(LINFO, format, false, v...)
}

func (l *Logger) Warn(v ...interface{}) {
	l.write(LWARN, "", false, v...)
}
func (l *Logger) Warnln(v ...interface{}) {
	l.write(LWARN, "", true, v...)
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.write(LWARN, format, false, v...)
}

func (l *Logger) Error(v ...interface{}) {
	l.write(LERROR, "", false, v...)
}
func (l *Logger) Errorln(v ...interface{}) {
	l.write(LERROR, "", true, v...)
}
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.write(LERROR, format, false, v...)
}

func (l *Logger) Fatal(v ...interface{}) {
	l.write(LFATAL, "", false, v...)
	os.Exit(1)
}
func (l *Logger) Fatalln(v ...interface{}) {
	l.write(LFATAL, "", true, v...)
	os.Exit(1)
}
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.write(LFATAL, format, false, v...)
	os.Exit(1)
}
