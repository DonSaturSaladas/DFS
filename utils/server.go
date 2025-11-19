package utils

import (
	"fmt"
	"log"
	"net"
)

func StartServer(port int, prefix string) (net.Listener, *log.Logger) {
	logger := CreateLogger(prefix)
	listener, err := ListenPort(port)
	if err != 0 {
		error_string := fmt.Sprintf("ERROR: No se pudo escuchar en el puerto %d\n", port)
		logger.Print(error_string)
		panic(error_string)
	}
	logger.Printf("INFO: Iniciado el servidor, escuchando el puerto %d\n", port)
	return listener, logger
}

func AcceptConnection(listener net.Listener, logger *log.Logger) net.Conn {
	conn, err := listener.Accept()
	if err != nil {
		error_string := "ERROR: No se pudo aceptar la conexion\n"
		logger.Print(error_string)
		panic(error_string)
	}
	return conn
}
