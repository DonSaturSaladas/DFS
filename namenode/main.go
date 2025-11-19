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
	listener, logger := utils.StartServer(port, "NAMENODE")
	for {
		conn := utils.AcceptConnection(listener, logger)
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	type Commands struct {
		Get string
		Put string
		Ls  string
	}
	commands := Commands{
		Get: "get",
		Put: "put",
		Ls:  "ls",
	}
	buffer := make([]byte, 256)
	commandsMap := map[string]func(string){
		commands.Get: getCommand,
		commands.Put: putCommand,
		commands.Ls:  lsCommand,
	}

	for {
		message := utils.ReadMessage(conn, buffer, logger)
		commandReaded := commandsMap[message.Command]
		if commandReaded == nil {
			conn.Write([]byte("El comando ingresado no es valido.\nIngrese uno de los siguientes comandos: get, put, ls\n"))
		} else {
			commandReaded(message.Parameters)
		}
	}
}

func getCommand(params string) {
	fmt.Printf("Get command con parametros %s\n", params)
}

func putCommand(params string) {
	fmt.Printf("Put command con parametros %s\n", params)
}

func lsCommand(params string) {
	fmt.Printf("Ls command con parametros %s\n", params)
}
