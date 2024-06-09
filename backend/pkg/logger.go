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

// color
const (
	CGreen        = "\033[32m"
	CWhite        = "\033[97m"
	CGray         = "\033[90m"
	CYellow       = "\033[93m"
	CRed          = "\033[91m"
	CWhiteFgRedBg = "\033[97;41m"
	CBlue         = "\033[94m"
	CMagenta      = "\033[35m"
	CCyan         = "\033[36m"
	CReset        = "\033[0m"
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

func GetLogLevel() int {
	var l int
	doWithLock(func() {
		l = logLevel
	})
	return l
}

type Logger struct {
	Prefix string
}

func NewLogger(prefix ...string) *Logger {
	l := &Logger{}
	if len(prefix) > 0 {
		l.Prefix = prefix[0]
	}
	return l
}

// default reset true
func (l *Logger) color(s string, clr string, reset ...bool) string {
	if !logColor || !logOutputIsTerm {
		return s
	}
	if len(reset) > 0 && !reset[0] {
		return clr + s
	}
	return clr + s + CReset

}
func (l *Logger) Green(s string, reset ...bool) string {
	return l.color(s, CGreen, reset...)
}

func (l *Logger) White(s string, reset ...bool) string {
	return l.color(s, CWhite, reset...)
}

func (l *Logger) Gray(s string, reset ...bool) string {
	return l.color(s, CGray, reset...)
}

func (l *Logger) Yellow(s string, reset ...bool) string {
	return l.color(s, CYellow, reset...)
}

func (l *Logger) Red(s string, reset ...bool) string {
	return l.color(s, CRed, reset...)
}

func (l *Logger) WhiteFgRedBg(s string, reset ...bool) string {
	return l.color(s, CWhiteFgRedBg, reset...)
}

func (l *Logger) Blue(s string, reset ...bool) string {
	return l.color(s, CBlue, reset...)
}

func (l *Logger) Magenta(s string, reset ...bool) string {
	return l.color(s, CMagenta, reset...)
}

func (l *Logger) Cyan(s string, reset ...bool) string {
	return l.color(s, CCyan, reset...)
}

func (l *Logger) GetResetString() string {
	return l.color("", CReset, false)
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
		green = CGreen
		white = CWhite
		gray = CGray
		yellow = CYellow
		red = CRed
		whiteFgRedBg = CWhiteFgRedBg
		blue = CBlue
		// magenta = CMagenta
		cyan = CCyan
		reset = CReset
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
	buf.WriteString(l.Prefix)
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
