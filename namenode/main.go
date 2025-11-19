package main

import (
	"dfs/utils"
	"log"
	"net"
)

var port int = 5000
var logger *log.Logger

func main() {
	logger = utils.CreateLogger("NAMENODE")
	listener := utils.StartServer(port, logger)
	for {
		conn := utils.AcceptConnection(listener, logger)
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	buffer := make([]byte, 256)
	commandsMap := map[string]func(net.Conn, string){
		"get": getCommand,
		"put": putCommand,
		"ls":  lsCommand,
	}

	for {
		message := utils.ReadMessage(conn, buffer, logger)
		commandReaded := commandsMap[message.Command]
		if commandReaded == nil {
			logger.Printf("INFO: El cliente ingreso un comando invalido: \"%s\"", message.Command)
			utils.SendMessage(conn, "El comando ingresado no es valido.\nIngrese uno de los siguientes comandos: get, put, ls\n")
		} else {
			commandReaded(conn, message.Parameters)
		}
	}
}
