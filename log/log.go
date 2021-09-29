package log

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

type Logger struct {
	mu     sync.Mutex
	logger *log.Logger
	file   string
	days   int
}

var logger = New("", 0, "")

func Default() *Logger {
	return logger
}

func New(file string, days int, prefix string) *Logger {
	w := os.Stdout
	if file != "" {
		w = LogFileWriter(FileNameToTodayName(file))
	}
	l := log.New(w, prefix, log.Ldate|log.Ltime|log.Lshortfile)
	return &Logger{logger: l, file: file, days: days}
}

func (l *Logger) GetStdLogger() *log.Logger {
	return l.logger
}

func (l *Logger) SetOutput(w io.Writer) {
	l.logger.SetOutput(w)
}

func (l *Logger) SetFlags(flag int) {
	l.logger.SetFlags(flag)
}

func (l *Logger) SetPrefix(prefix string) {
	l.logger.SetPrefix(prefix)
}

func (l *Logger) Println(v ...interface{}) {
	l.RotateLogFile()
	l.logger.Println(v...)
}

func (l *Logger) Fatalln(v ...interface{}) {
	l.RotateLogFile()
	l.logger.Fatalln(v...)
}

func (l *Logger) Panicln(v ...interface{}) {
	l.RotateLogFile()
	l.logger.Panicln(v...)
}

func (l *Logger) Panicf(format string, v ...interface{}) {
	l.RotateLogFile()
	l.logger.Panicf(format, v...)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.RotateLogFile()
	l.logger.Fatalf(format, v...)
}

func (l *Logger) Printf(format string, v ...interface{}) {
	l.RotateLogFile()
	l.logger.Printf(format, v...)
}

func (l *Logger) RotateLogFile() {
	if l.file != "" {
		l.mu.Lock()
		defer l.mu.Unlock()
		file := FileNameToTodayName(l.file)
		if _, err := os.Stat(file); os.IsNotExist(err) {
			l.logger.Writer().(*os.File).Close()
			l.SetOutput(LogFileWriter(file))
		}
		//delete files
		delFiles := l.DeleteFiles()
		for _, f := range delFiles {
			os.Remove(f)
		}
	}
}

func (l *Logger) DeleteFiles() []string {
	var delfiles = []string{}
	if l.days <= 0 {
		return delfiles
	}

	dir, ext, _, basename := GetFileInfo(l.file)
	glob := dir + "/" + basename + "-2[0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]" + ext
	files, err := filepath.Glob(glob)
	if err == nil {
		if l.days < len(files) {
			sort.Strings(files)
			delfiles = files[0 : len(files)-l.days]
		}
	}
	return delfiles
}

func LogFileWriter(filepath string) *os.File {
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	return file
}

func GetFileInfo(file string) (string, string, string, string) {
	dir := filepath.Dir(file)
	ext := filepath.Ext(file)
	filename := filepath.Base(file)
	basename := filename

	dotpos := strings.LastIndex(filename, ".")
	if dotpos > 0 {
		basename = filename[0:dotpos]
	}
	return dir, ext, filename, basename
}

func FileNameToDateName(file string, time time.Time) string {
	dir, ext, _, basename := GetFileInfo(file)
	return dir + "/" + basename + time.Format("-2006-01-02") + ext
}

func FileNameToTodayName(file string) string {
	return FileNameToDateName(file, time.Now())
}

func Println(v ...interface{}) {
	logger.Println(v...)
}

func Fatalln(v ...interface{}) {
	logger.Fatalln(v...)
}

func Panicln(v ...interface{}) {
	logger.Panicln(v...)
}

func Panicf(format string, v ...interface{}) {
	logger.Panicf(format, v...)
}

func Fatalf(format string, v ...interface{}) {
	logger.Fatalf(format, v...)
}

func Printf(format string, v ...interface{}) {
	logger.Printf(format, v...)
}
