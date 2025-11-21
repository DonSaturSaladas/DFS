package main

import (
	"dfs/utils"
	"fmt"
	"log"
	"net"
)

var logger *log.Logger
var ip string //LOCALMENTE (borrar cuando este finalizado)

func main() {
	listener := startDatanode()
	for {
		conn := utils.AcceptConnection(listener, logger)
		go handleConnection(conn)
	}
}

func startDatanode() net.Listener {
	fmt.Print("Ingresar la ip del nodo: ")
	fmt.Scan(&ip)
	logger = utils.CreateLogger("NAMENODE")
	listener := utils.StartServer(5000, logger)
	return listener
}

func handleConnection(conn net.Conn) {
	commandsMap := map[string]func(utils.Message){
		"store": storeCommand,
		"read":  readCommand,
	}

	for {
		message := utils.ReadMessage(conn, logger)
		commandReaded := commandsMap[message.Command]
		if commandReaded == nil {
			logger.Printf("INFO: El cliente ingreso un comando invalido: \"%s\"", message.Command)
			utils.SendMessage(conn, "El comando ingresado no es valido.\nIngrese uno de los siguientes comandos: store, read\n")
		} else {
			commandReaded(message)
		}
	}
}
