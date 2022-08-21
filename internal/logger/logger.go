package logger

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	logger *logrus.Logger
}

func New(level string) (*Logger, error) {
	logg := logrus.New()

	logrusLevel, err := logrus.ParseLevel(level)
	if err != nil {
		return nil, err
	}

	logg.SetLevel(logrusLevel)

	return &Logger{logg}, nil
}

func (l Logger) Error(msg string) {
	l.logger.Error(msg)
}

func (l Logger) Warning(msg string) {
	l.logger.Warning(msg)
}

func (l Logger) Info(msg string) {
	l.logger.Info(msg)
}

func (l Logger) Debug(msg string) {
	l.logger.Debug(msg)
}

func (l *Logger) LogHTTPRequest(r *http.Request, statusCode int, requestDuration time.Duration) {
	l.logger.Infof(
		"%s [%s] %s %s %s %d %s %q",
		r.RemoteAddr,
		time.Now().Format(time.RFC1123Z),
		r.Method,
		r.RequestURI,
		r.Proto,
		statusCode,
		requestDuration,
		r.UserAgent(),
	)
}

func (l *Logger) LogGRPCRequest(r interface{}, method string, requestDuration time.Duration) {
	l.logger.Infof(
		"[%s] %s %s %s",
		time.Now().Format(time.RFC1123Z),
		method,
		requestDuration,
		r,
	)
}
