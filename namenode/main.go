package main

import (
	"dfs/utils"
	"encoding/json"
	"log"
	"net"
	"os"
)

var logger *log.Logger
var systemInfo SystemInfo

func main() {
	listener := startNamenode()
	for {
		conn := utils.AcceptConnection(listener, logger)
		go handleConnection(conn)
	}
}

func startNamenode() net.Listener {
	logger = utils.CreateLogger("NAMENODE")
	listener := utils.StartServer(5000, logger)
	parseSystemInfo()
	return listener
}

func parseSystemInfo() {
	byteValue, err := os.ReadFile("data/system_info.json")
	if err != nil {
		logger.Print("ERROR: No se pudo abrir el archivo de informacion del sistema.")
		panic(err)
	}
	json.Unmarshal(byteValue, &systemInfo)
}

func handleConnection(conn net.Conn) {
	commandsMap := map[string]func(utils.Message){
		"get": getCommand,
		"put": putCommand,
		"ls":  lsCommand,
	}

	for {
		message := utils.ReadMessage(conn, logger)
		commandReaded := commandsMap[message.Command]
		if commandReaded == nil {
			logger.Printf("INFO: El cliente ingreso un comando invalido: \"%s\"", message.Command)
			utils.SendMessage(conn, "El comando ingresado no es valido.\nIngrese uno de los siguientes comandos: get, put, ls\n")
		} else {
			commandReaded(message)
		}
	}
}
