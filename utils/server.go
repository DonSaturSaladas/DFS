package utils

import (
	"fmt"
	"log"
	"net"
)

func StartServer(port int, logger *log.Logger) net.Listener {
	listener, err := ListenPort(port)
	if err != 0 {
		errorString := fmt.Sprintf("ERROR: No se pudo escuchar en el puerto %d\n", port)
		logger.Print(errorString)
		panic(errorString)
	}
	logger.Printf("INFO: Iniciado el servidor, escuchando el puerto %d\n", port)
	return listener
}

func AcceptConnection(listener net.Listener, logger *log.Logger) net.Conn {
	conn, err := listener.Accept()
	if err != nil {
		errorString := "ERROR: No se pudo aceptar la conexion\n"
		logger.Print(errorString)
		panic(errorString)
	}
	logger.Print("INFO: Se establecio conexion con un cliente\n")
	return conn
}
