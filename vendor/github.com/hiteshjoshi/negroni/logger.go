package negroni

import (
	"bytes"

	"github.com/fatih/color"

	"log"
	"net/http"
	"os"
	"text/template"
	"time"
)

// LoggerEntry is the structure
// passed to the template.
type LoggerEntry struct {
	AppName   string
	StartTime string
	Status    int
	Duration  time.Duration
	Hostname  string
	Method    string
	Path      string
}

// LoggerDefaultFormat is the format
// logged used by the default Logger instance.
var LoggerDefaultFormat = "{{.AppName}} > {{.StartTime}} | {{.Status}} | \t {{.Duration}} | {{.Hostname}} | {{.Method}} {{.Path}} \n"

// LoggerDefaultDateFormat is the
// format used for date by the
// default Logger instance.
var LoggerDefaultDateFormat = time.RFC3339

// ALogger interface
type ALogger interface {
	Print(v ...interface{})
	Printf(format string, v ...interface{})
}

// Logger is a middleware handler that logs the request as it goes in and the response as it goes out.
type Logger struct {
	AppName string
	// ALogger implements just enough log.Logger interface to be compatible with other implementations
	ALogger
	dateFormat string
	template   *template.Template
}

// NewLogger returns a new Logger instance
func NewLogger() *Logger {
	logger := &Logger{ALogger: log.New(os.Stdout, "[negroni] ", 0), dateFormat: LoggerDefaultDateFormat}
	logger.SetFormat(LoggerDefaultDateFormat)
	return logger
}

func (l *Logger) SetAppName(name string) {
	l.AppName = name
}
func (l *Logger) SetFormat(format string) {
	l.template = template.Must(template.New("negroni_parser").Parse(format))
}

func (l *Logger) SetDateFormat(format string) {
	l.dateFormat = format
}

func (l *Logger) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()

	next(rw, r)

	res := rw.(ResponseWriter)
	log := LoggerEntry{
		AppName:   l.AppName,
		StartTime: start.Format(l.dateFormat),
		Status:    res.Status(),
		Duration:  time.Since(start),
		Hostname:  r.Host,
		Method:    r.Method,
		Path:      r.URL.Path,
	}

	buff := &bytes.Buffer{}
	l.template.Execute(buff, log)

	d := color.New(color.FgCyan, color.Bold)
	d.Printf(buff.String())

}
