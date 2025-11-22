package main

import (
	"dfs/utils"
)

func lsCommand(data utils.Message) {
	metadataMutex.RLock()
	filesList := readMetadataFiles(metadata)
	metadataMutex.RUnlock()

	message := getLsMessage(data, filesList)

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
		command = "nofiles"
	} else {
		command = "files"
	}
	message := utils.Message{
		Connection: data.Connection,
		Command:    command,
		Parameters: fileList,
	}
	return message
}
