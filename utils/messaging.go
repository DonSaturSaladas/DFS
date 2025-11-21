package utils

import (
	"log"
	"net"
	"strings"
)

type Message struct {
	Command    string
	Parameters []string
	Connection net.Conn
}

func SendMessage(conn net.Conn, message string) {
	byteMessage := []byte(message)
	conn.Write(byteMessage)
}

func ReadMessage(conn net.Conn, logger *log.Logger) Message {
	buffer := make([]byte, 1024)
	readedBytes := readConnection(conn, buffer)
	checkReadingError(readedBytes, logger)
	messageString := string(buffer[:readedBytes])
	message := parseMessageData(messageString)
	message.Connection = conn
	return message
}

func checkReadingError(readedBytes int, logger *log.Logger) {
	if readedBytes < 0 {
		errorString := "ERROR: No se pudo leer en la conexion\n"
		logger.Print(errorString)
		panic(errorString)
	}
}

func readConnection(conn net.Conn, buffer []byte) int {
	readedBytes, err := conn.Read(buffer)
	if err != nil {
		readedBytes = -1
	}
	return readedBytes
}

func parseMessageData(message string) Message {
	message = strings.TrimSuffix(strings.TrimSuffix(message, "\n"), "\r")
	messageSlice := strings.Split(message, " ")
	var params []string
	if len(messageSlice) > 1 {
		params = messageSlice[1:]
	}
	msg := Message{
		Command:    messageSlice[0],
		Parameters: params,
	}
	return msg
}
