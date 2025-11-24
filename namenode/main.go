package main

import (
	"dfs/utils"
	"encoding/json"
	"log"
	"net"
	"os"
	"sync"
)

var logger *log.Logger

var systemInfo SystemInfo
var systemInfoMutex sync.RWMutex

var metadata map[string][]Block
var metadataMutex sync.RWMutex

func main() {
	listener := startNamenode()
	defer listener.Close()
	for {
		conn := utils.AcceptConnection(listener, logger)
		go handleConnection(conn)
	}
}

func startNamenode() net.Listener {
	logger = utils.CreateLogger("NAMENODE")
	listener := utils.StartServer(5000, logger)
	parseSystemInfo()
	parseMetadata()
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

func parseMetadata() {
	byteValue, err := os.ReadFile("data/metadata.json")
	if err != nil {
		logger.Print("ERROR: No se pudo abrir el archivo de metadata.")
		panic(err)
	}
	json.Unmarshal(byteValue, &metadata)
}

func handleConnection(conn net.Conn) {
	commandsMap := map[string]func(utils.Message){
		"get": getCommand,
		"put": putCommand,
		"ls":  lsCommand,
		"rm":  rmCommand,
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
