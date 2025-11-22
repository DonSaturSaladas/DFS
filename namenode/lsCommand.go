package main

import (
	"dfs/utils"
)

func lsCommand(data utils.Message) {
	metadataMutex.RLock()
	response := readMetadataFiles(metadata)
	metadataMutex.RUnlock()

	utils.SendMessage(data.Connection, response)
}

func readMetadataFiles(metadata map[string][]Block) string {
	response := "["
	for fileName := range metadata {
		response = response + fileName + ","
	}
	return response[:len(response)-1] + "]"
}
