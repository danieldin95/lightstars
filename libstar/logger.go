package libstar

import (
	"container/list"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

const (
	PRINT = 0x00
	DEUBG = 0x01
	INFO  = 0x02
	WARN  = 0x03
	ERROR = 0x04
	FATAL = 0xff
)

type Logger struct {
	Level    int
	FileName string
	FileLog  *log.Logger
	Lock     sync.Mutex
	Errors   *list.List
}

type Message struct {
	Level   string `json:"level"`
	Date    string `json:"date"`
	Message string `json:"message"`
}

func (l *Logger) Debug(format string, v ...interface{}) {
	if DEUBG >= l.Level {
		log.Printf("DEBUG "+format, v...)
	}
}

func (l *Logger) Info(format string, v ...interface{}) {
	if INFO >= l.Level {
		log.Printf("INFO "+format, v...)
	}
	l.SaveError("INFO", format, v...)
}

func (l *Logger) Warn(format string, v ...interface{}) {
	if WARN >= l.Level {
		log.Printf("WARN"+format, v...)
	}
	l.SaveError("WARN", format, v...)
}

func (l *Logger) Error(format string, v ...interface{}) {
	if ERROR >= l.Level {
		log.Printf("ERROR "+format, v...)
	}
	l.SaveError("ERROR", format, v...)
}

func (l *Logger) Fatal(format string, v ...interface{}) {
	if FATAL >= l.Level {
		log.Printf("FATAL "+format, v...)
	}
	l.SaveError("FATAL", format, v...)
}

func (l *Logger) Print(format string, v ...interface{}) {
	if PRINT >= l.Level {
		log.Printf("PRINT"+format, v...)
	}
}

func (l *Logger) SaveError(level string, format string, v ...interface{}) {
	m := fmt.Sprintf(format, v...)
	if l.FileLog != nil {
		l.FileLog.Println(level+" "+m)
	}

	l.Lock.Lock()
	defer l.Lock.Unlock()
	if l.Errors.Len() >= 1024 {
		if e := l.Errors.Back(); e != nil {
			l.Errors.Remove(e)
		}
	}
	yy, mm, dd := time.Now().Date()
	hh, mn, se := time.Now().Clock()
	ele := &Message{
		Level:   level,
		Date:    fmt.Sprintf("%d/%d/%d %d:%d:%d", yy, mm, dd, hh, mn, se),
		Message: m,
	}
	l.Errors.PushBack(ele)
}

func (l *Logger) List() <-chan *Message {
	c := make(chan *Message, 128)
	go func() {
		l.Lock.Lock()
		defer l.Lock.Unlock()
		for ele := l.Errors.Back(); ele != nil; ele = ele.Prev() {
			c <- ele.Value.(*Message)
		}
		c <- nil // Finish channel by nil.
	}()
	return c
}

var Log = Logger{
	Level:    INFO,
	FileName: ".log.error",
	Errors:   list.New(),
}

func Print(format string, v ...interface{}) {
	Log.Print(format, v...)
}

func Error(format string, v ...interface{}) {
	Log.Error(format, v...)
}

func Debug(format string, v ...interface{}) {
	Log.Debug(format, v...)
}

func Info(format string, v ...interface{}) {
	Log.Info(format, v...)
}

func Warn(format string, v ...interface{}) {
	Log.Warn(format, v...)
}

func Fatal(format string, v ...interface{}) {
	Log.Fatal(format, v...)
}

func Init(file string, level int) {
	SetLog(level)
	Log.FileName = file
	if Log.FileName != "" {
		logFile, err := os.Create(Log.FileName)
		if err == nil {
			Log.FileLog = log.New(logFile, "", log.LstdFlags)
		} else {
			Warn("logger.Init: %s", err)
		}
	}
}

func SetLog(level int) {
	Log.Level = level
}

func Close() {
	//TODO
}

func Catch(name string) {
	if err := recover(); err != nil {
		Fatal("%s Panic: ===%s===", name, err)
		Fatal("%s Stack: ===%s===", name, debug.Stack())
	}
}
