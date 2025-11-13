package main

import (
	"dfs/utils"
	"fmt"
	"net"
)

var port int = 5000

func main() {
	logger := utils.CreateLogger("NAMENODE")
	listener, err := utils.ListenPort(port)
	if err != 0 {
		logger.Printf("ERROR: No se pudo escuchar en el puerto %d\n", port)
	}
	logger.Printf("INFO: Iniciado el servidor, escuchando el puerto %d\n", port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go handleConnection(conn)
	}
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
