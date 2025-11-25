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
	logger.Printf("INFO: Cargando los datos del archivo \"system_info.json\".\n")
	byteValue, err := os.ReadFile("data/system_info.json")
	if err != nil {
		logger.Print("ERROR: No se pudo abrir el archivo de informacion del sistema.")
		panic(err)
	}
	json.Unmarshal(byteValue, &systemInfo)
	logger.Printf("INFO: Datos del archivo \"system_info.json\" cargados.\n")
}

func parseMetadata() {
	logger.Printf("INFO: Cargando los datos del archivo \"metadata.json\".\n")
	byteValue, err := os.ReadFile("data/metadata.json")
	if err != nil {
		logger.Print("ERROR: No se pudo abrir el archivo de metadata.")
		panic(err)
	}
	json.Unmarshal(byteValue, &metadata)
	logger.Printf("INFO: Datos del archivo \"metadata.json\" cargados.\n")
}

func handleConnection(conn net.Conn) {
	commandsMap := map[string]func(utils.Message){
		"get": getCommand,
		"put": putCommand,
		"ls":  lsCommand,
		"rm":  rmCommand,
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
