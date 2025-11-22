package main

import (
	"dfs/utils"
)

func getCommand(data utils.Message) {
	fileName := data.Parameters[0]
	metadataMutex.RLock()
	fileMetadata := metadata[fileName]
	addresses := readFileMetadata(fileMetadata)
	metadataMutex.RUnlock()
	message := utils.Message{
		Connection: data.Connection,
		Command:    "addresses",
		Parameters: addresses,
	}
	message.Send()
}

func readFileMetadata(fileMetadata []Block) []string {
	addresses := make([]string, len(fileMetadata))
	for index, block := range fileMetadata {
		addresses[index] = block.Adress
	}
	return addresses
}
