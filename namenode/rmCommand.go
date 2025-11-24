package main

import (
	"dfs/utils"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"slices"
)

func rmCommand(data utils.Message) {
	fileName := data.Parameters[0]
	logger.Printf("INFO: Buscando en el archivo de metadata los nodos que contienen al archivo \"%s\".\n", fileName)
	addresses := getFileAddresses(fileName)
	if len(addresses) == 0 {
		noBlocksToRemove(data.Connection)
		return
	}
	sendBlockAddressesToRemove(data.Connection, addresses)
	response := utils.ReadMessage(data.Connection, logger)
	if response.Command == "confirm" {
		fmt.Print("Confirmado\n")
		removeFileData(fileName)
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
	metadataMutex.Lock()
	fileMetadata := metadata[fileName]
	delete(metadata, fileName)
	metadataContent, _ := json.MarshalIndent(metadata, "", "\t")
	file, _ := os.Create("data/metadata.json")
	file.Write(metadataContent)
	metadataMutex.Unlock()
	return fileMetadata
}

func removeFileFromSystemInfo(fileMetadata []Block) {
	nodeUsageInFile := getNodeUsageInFile(fileMetadata)
	systemInfoMutex.Lock()
	for nodeAddress, usageInFile := range nodeUsageInFile {
		removeNodeUsage(nodeAddress, usageInFile)
	}
	systemInfoMutex.Unlock()
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
