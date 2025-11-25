package utils

import (
	"fmt"
	"io"
	"log"
	"os"
)

func CreateLogger(prefix string) *log.Logger {
	if prefix == "" {
		panic("Ingrese un prefijo\n")
	}
	logFile := OpenFile("data/logs.log")
	logger := log.New(io.MultiWriter(os.Stdout, logFile), fmt.Sprintf("[%s] ", prefix), log.Ldate|log.Ltime)
	return logger
}

func OpenFile(path string) *os.File {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	return file
}
