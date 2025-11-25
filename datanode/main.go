package main

import (
	"dfs/utils"
	"log"
	"net"
)

var logger *log.Logger

func main() {
	listener := startDatanode()
	defer listener.Close()
	for {
		conn := utils.AcceptConnection(listener, logger)
		go handleConnection(conn)
	}
}

func startDatanode() net.Listener {
	logger = utils.CreateLogger("DATANODE")
	listener := utils.StartServer(5000, logger)
	return listener
}

func handleConnection(conn net.Conn) {
	commandsMap := map[string]func(utils.Message){
		"store": storeCommand,
		"read":  readCommand,
		"rm":    rmCommand,
	}

	for {
		logger.Print("INFO: Esperando un comando.\n")
		message := utils.ReadMessage(conn, logger)
		commandReaded := commandsMap[message.Command]
		if message.Command == "connection_ended" {
			logger.Print("INFO: Conexion finalizada, terminando hilo.\n")
			return
		} else {
			logger.Printf("INFO: Comando recibido -> %s.\n", message.Command)
			commandReaded(message)
		}
	}
}
