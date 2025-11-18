package main

import (
	"dfs/utils"
	"fmt"
	"log"
	"net"
)

var port int = 5000
var logger *log.Logger

func main() {
	listener := startServer()
	for {
		conn := acceptConnection(listener)
		go handleConnection(conn)
	}
}

func startServer() net.Listener {
	logger = utils.CreateLogger("NAMENODE")
	listener, err := utils.ListenPort(port)
	if err != 0 {
		error_string := fmt.Sprintf("ERROR: No se pudo escuchar en el puerto %d\n", port)
		logger.Print(error_string)
		panic(error_string)
	}
	logger.Printf("INFO: Iniciado el servidor, escuchando el puerto %d\n", port)
	return listener
}

func acceptConnection(listener net.Listener) net.Conn {
	conn, err := listener.Accept()
	if err != nil {
		error_string := "ERROR: No se pudo aceptar la conexion\n"
		logger.Print(error_string)
		panic(error_string)
	}
	return conn
}

func handleConnection(conn net.Conn) {
	buffer := make([]byte, 256)
	for {
		_, err := conn.Read(buffer)
		fmt.Printf("Ingresado algo\n")
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s", buffer)
	}
}
