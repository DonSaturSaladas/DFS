package main

import (
	"dfs/utils"
	"fmt"
	"log"
	"net"
	"strconv"
)

var logger *log.Logger
var puerto string //LOCALMENTE (borrar cuando este finalizado)

func main() {
	listener := startDatanode()
	for {
		conn := utils.AcceptConnection(listener, logger)
		go handleConnection(conn)
	}
}

func startDatanode() net.Listener {
	fmt.Print("Ingresar el puerto del nodo: ")
	fmt.Scan(&puerto)
	port, _ := strconv.Atoi(puerto)
	logger = utils.CreateLogger("DATANODE")
	listener := utils.StartServer(port, logger)
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
		if message.Command == "connection_ended" {
			return
		} else {
			commandReaded(message)
		}
	}
}
