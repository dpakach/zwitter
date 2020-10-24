package log

import (
	"log"
	"os"
)

const colorReset = "\033[0m"
const colorRed = "\033[31m"
const colorGreen = "\033[32m"
const colorYellow = "\033[33m"

type ZwitLogger struct {
	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
}

func New() *ZwitLogger {
	return &ZwitLogger{
		infoLogger:    log.New(os.Stdout, string(colorGreen)+"[INFO]: "+string(colorReset), log.LstdFlags),
		warningLogger: log.New(os.Stdout, string(colorYellow)+"[WARN]: "+string(colorReset), log.LstdFlags),
		errorLogger:   log.New(os.Stdout, string(colorRed)+"[ERROR]: "+string(colorReset), log.LstdFlags),
	}
}

func (zl *ZwitLogger) Info(text string) {
	zl.infoLogger.Println(text)
}

func (zl *ZwitLogger) Warn(text string) {
	zl.warningLogger.Println(text)
}

func (zl *ZwitLogger) Error(text string) {
	zl.errorLogger.Println(text)
}

func (zl *ZwitLogger) Infof(text string, args ...interface{}) {
	zl.infoLogger.Printf(text, args...)
}

func (zl *ZwitLogger) Warnf(text string, args ...interface{}) {
	zl.warningLogger.Printf(text, args...)
}

func (zl *ZwitLogger) Errorf(text string, args ...interface{}) {
	zl.errorLogger.Printf(text, args...)
}
