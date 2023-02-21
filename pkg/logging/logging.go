package logging

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"runtime"
	"time"
)

type writerHook struct {
	Writer    []io.Writer
	LogLevels []logrus.Level
}

func (hook *writerHook) Fire(entry *logrus.Entry) error {
	line, err := entry.Bytes()
	if err != nil {
		return err
	}
	for _, w := range hook.Writer {
		w.Write(line)
	}
	return err
}

func (hook *writerHook) Levels() []logrus.Level {
	return hook.LogLevels
}

var e *logrus.Entry

type Logger struct {
	*logrus.Entry
}

func GetLogger() Logger {
	return Logger{e}
}

func (l *Logger) GetLoggerWithFiled(k string, v interface{}) Logger {
	return Logger{l.WithField(k, v)}
}

func init() {
	log := logrus.New()
	log.SetReportCaller(true)
	log.Formatter = &logrus.JSONFormatter{
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			fileName := path.Base(frame.File)
			return fmt.Sprintf("%s()", frame.Function), fmt.Sprintf("%s: %d", fileName, frame.Line)
		},
		DataKey:           time.Now().String(),
		TimestampFormat:   time.RFC3339,
		DisableHTMLEscape: false,
		DisableTimestamp:  true,
	}
	err := os.MkdirAll("logs", 0644)
	if err != nil {
		panic(err)
	}
	allFile, err := os.OpenFile("logs/all.log", os.O_CREATE|os.O_APPEND, 0640)
	if err != nil {
		panic(err)
	}
	log.SetOutput(io.Discard)

	log.AddHook(&writerHook{
		Writer:    []io.Writer{allFile},
		LogLevels: logrus.AllLevels,
	})
	log.SetLevel(logrus.FatalLevel)
	log.SetLevel(logrus.ErrorLevel)
	log.SetLevel(logrus.TraceLevel)
	e = logrus.NewEntry(log)
}
