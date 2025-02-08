package logger

import (
	"log"
	"os"
	"time"
)

// Loggers struct contains two loggers: one for info logs and one for error logs.
type Loggers struct {
	LogInfo  *log.Logger
	LogError *log.Logger
}

// SetLoggers initializes and returns a Loggers struct with info and error loggers.
func SetLoggers() *Loggers {
	serverLogFile := buildLogFiles()
	log.SetFlags(log.Ldate | log.Ltime)
	loggers := new(Loggers)
	loggers.LogInfo = log.New(serverLogFile, "Info: ", log.Ldate|log.Ltime)
	loggers.LogError = log.New(serverLogFile, "Error: ", log.Ldate|log.Ltime)
	return loggers
}

// buildLogFiles creates a directory for log files if it does not exist.
func buildLogFiles() (serverLogFile *os.File) {
	time := time.Now().Format(time.DateOnly)
	dir := "logs/"
	err := os.MkdirAll(dir, 0o777)
	CheckError(err)
	serverLogFile, err = os.OpenFile(dir+time+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	CheckError(err)
	return serverLogFile
}

// CheckError is a helper function to check if an error occurred.
func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
