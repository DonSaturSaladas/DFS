package utils

import (
	"fmt"
	"log"
	"os"
)

var baseDir string = "/home/maxi/uns/soyd/labos/DFS"

func CreateLogger(prefix string) *log.Logger {
	if prefix == "" {
		panic("Ingrese un prefijo\n")
	}
	logFile := openLogFile()
	logger := log.New(logFile, fmt.Sprintf("[%s] ", prefix), log.Ldate|log.Ltime)
	return logger
}

func openLogFile() *os.File {
	logsPath := fmt.Sprintf("%s/data/logs.log", baseDir)
	file, err := os.OpenFile(logsPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	return file
}
