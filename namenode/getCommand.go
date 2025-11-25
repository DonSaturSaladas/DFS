package main

import (
	"dfs/utils"
)

func getCommand(data utils.Message) {
	fileName := data.Parameters[0]
	logger.Print("INFO: Solicitando mutex para leer la informacion de \"metadata.json\".\n")
	metadataMutex.RLock()
	logger.Print("INFO: Obtenido mutex de \"metadata.json\".\n")
	fileMetadata := metadata[fileName]
	logger.Printf("INFO: Buscando las direcciones del archivo \"%s\".\n", fileName)
	addresses := readFileMetadata(fileMetadata)
	metadataMutex.RUnlock()
	logger.Print("INFO: Liberado mutex de \"metadata.json\".\n")
	message := utils.Message{
		Connection: data.Connection,
		Command:    "addresses",
		Parameters: addresses,
	}
	logger.Printf("INFO: Enviando las direcciones del archivo \"%s\".\n", fileName)
	message.Send()
}

func readFileMetadata(fileMetadata []Block) []string {
	addresses := make([]string, len(fileMetadata))
	for index, block := range fileMetadata {
		addresses[index] = block.Adress
	}
	return addresses
}
