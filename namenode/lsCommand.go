package main

import (
	"dfs/utils"
)

func lsCommand(data utils.Message) {
	logger.Print("INFO: Solicitando mutex para leer la informacion de \"metadata.json\".\n")
	metadataMutex.RLock()
	logger.Print("INFO: Obtenido mutex de \"metadata.json\".\n")
	logger.Print("INFO: Obteniendo los archivos disponibles en el DFS.\n")
	filesList := readMetadataFiles(metadata)
	metadataMutex.RUnlock()
	logger.Print("INFO: Liberado mutex de \"metadata.json\".\n")

	message := getLsMessage(data, filesList)
	logger.Print("INFO: Enviando los archivos existentes en el DFS.\n")
	message.Send()
}

func readMetadataFiles(metadata map[string][]Block) []string {
	fileNames := make([]string, len(metadata))
	i := 0
	for fileName := range metadata {
		fileNames[i] = fileName
		i++
	}
	return fileNames
}

func getLsMessage(data utils.Message, fileList []string) utils.Message {
	var command string
	if len(fileList) == 0 {
		logger.Print("INFO: El DFS no tiene ningun archivo.\n")
		command = "nofiles"
	} else {
		logger.Print("INFO: Existen archivos en el DFS.\n")
		command = "files"
	}
	message := utils.Message{
		Connection: data.Connection,
		Command:    command,
		Parameters: fileList,
	}
	return message
}
