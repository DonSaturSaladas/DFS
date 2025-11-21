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
func SendData(conn net.Conn, data []byte) {
	conn.Write(data)
}

func ReadMessage(conn net.Conn, logger *log.Logger) Message {
	buffer := make([]byte, 256)
	readedBytes := readConnection(conn, buffer)
	checkReadingError(readedBytes, logger)
	messageString := string(buffer[:readedBytes])
	message := ParseMessageData(messageString)
	message.Connection = conn
	return message
}

func ReadData(conn net.Conn, buffer []byte, logger *log.Logger) int {
	readedBytes := readConnection(conn, buffer)
	checkReadingError(readedBytes, logger)
	return readedBytes
}

func readConnection(conn net.Conn, buffer []byte) int {
	readedBytes, err := conn.Read(buffer)
	if err != nil {
		readedBytes = -1
	}
	return readedBytes
}

func checkReadingError(readedBytes int, logger *log.Logger) {
	if readedBytes < 0 {
		errorString := "ERROR: No se pudo leer en la conexion\n"
		logger.Print(errorString)
		panic(errorString)
	}
}

func ParseMessageData(message string) Message {
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
