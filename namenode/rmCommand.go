package main

import (
	"dfs/utils"
	"encoding/json"
	"net"
	"os"
	"slices"
)

func rmCommand(data utils.Message) {
	fileName := data.Parameters[0]
	logger.Printf("INFO: Buscando en el archivo de metadata los nodos que contienen al archivo \"%s\".\n", fileName)
	addresses := getFileAddresses(fileName)
	if len(addresses) == 0 {
		logger.Printf("INFO: El archivo \"%s\" no existe en el DFS.\n", fileName)
		noBlocksToRemove(data.Connection)
		return
	}
	logger.Printf("INFO: El archivo \"%s\" esta en el DFS.\n", fileName)
	logger.Printf("INFO: Enviando las direcciones de los nodos que tienen bloques a eliminar del archivo \"%s\".\n", fileName)
	sendBlockAddressesToRemove(data.Connection, addresses)
	logger.Printf("INFO: Esperando confirmacion de eliminacion de los bloques del archivo \"%s\".\n", fileName)
	response := utils.ReadMessage(data.Connection, logger)
	if response.Command == "confirm" {
		logger.Printf("INFO: Se recibio confirmacion de la eliminacion de los bloques del archivo \"%s\".\n", fileName)
		logger.Printf("INFO: Eliminando la informacion del sistema del archivo \"%s\".\n", fileName)
		removeFileData(fileName)
		logger.Printf("INFO: Informacion del archivo \"%s\" eliminada.\n", fileName)
	}

}

func getFileAddresses(fileName string) []string {
	metadataMutex.RLock()
	fileBlocks := metadata[fileName]
	metadataMutex.RUnlock()

	nodeAddresses := []string{}
	for _, block := range fileBlocks {
		if !slices.Contains(nodeAddresses, block.Adress) {
			nodeAddresses = append(nodeAddresses, block.Adress)
		}
	}
	return nodeAddresses
}

func noBlocksToRemove(conn net.Conn) {
	utils.Message{
		Connection: conn,
		Command:    "nofile",
	}.Send()
}

func sendBlockAddressesToRemove(conn net.Conn, addresses []string) {
	utils.Message{
		Connection: conn,
		Command:    "addresses",
		Parameters: addresses,
	}.Send()
}

func removeFileData(fileName string) {
	fileMetadata := removeFileFromMetadata(fileName)
	removeFileFromSystemInfo(fileMetadata)
}

func removeFileFromMetadata(fileName string) []Block {
	logger.Print("INFO: Solicitando mutex para escribir en el archivo \"metadata.json\".\n")
	metadataMutex.Lock()
	fileMetadata := metadata[fileName]
	delete(metadata, fileName)
	metadataContent, _ := json.MarshalIndent(metadata, "", "\t")
	file, _ := os.Create("data/metadata.json")
	file.Write(metadataContent)
	metadataMutex.Unlock()
	logger.Print("INFO: Liberado mutex de \"metadata.json\".\n")
	return fileMetadata
}

func removeFileFromSystemInfo(fileMetadata []Block) {
	nodeUsageInFile := getNodeUsageInFile(fileMetadata)
	logger.Print("INFO: Solicitando mutex para escribir en el archivo \"system_info.json\".\n")
	systemInfoMutex.Lock()
	for nodeAddress, usageInFile := range nodeUsageInFile {
		removeNodeUsage(nodeAddress, usageInFile)
	}
	systemInfoMutex.Unlock()
	logger.Print("INFO: Liberado mutex de \"system_info.json\".\n")
}

func getNodeUsageInFile(fileMetadata []Block) map[string]int {
	nodeUsageInFile := map[string]int{}
	for _, block := range fileMetadata {
		nodeUsageInFile[block.Adress]++
	}
	return nodeUsageInFile
}

func removeNodeUsage(nodeAddress string, usage int) {
	nodesBlockUsage := make([]BlockInfo, len(systemInfo.BlockUsage))
	copy(nodesBlockUsage, systemInfo.BlockUsage)

	for nodeIndex := range nodesBlockUsage {
		if nodesBlockUsage[nodeIndex].Address == nodeAddress {
			nodesBlockUsage[nodeIndex].Usage -= usage
		}
	}
	saveBlockUsageModifications(nodesBlockUsage)
}
