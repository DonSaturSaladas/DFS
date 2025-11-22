package main

import (
	"dfs/utils"
	"fmt"
)

func getCommand(data utils.Message) {
	if checkEmptyParams(data) {
		return
	}
	fileName := data.Parameters[0]
	metadataMutex.RLock()
	fileMetadata := metadata[fileName]
	response := readFileMetadata(fileMetadata)
	metadataMutex.RUnlock()
	utils.SendMessage(data.Connection, "addresses "+response)
}

func checkEmptyParams(data utils.Message) bool {
	if len(data.Parameters) == 0 {
		logger.Printf("ERROR: El cliente ingreso el comando %s sin parametros.\n", data.Command)
		utils.SendMessage(data.Connection, fmt.Sprintf("Ingrese los parametros luego del comando \"%s\"\n", data.Command))
		return true
	}
	return false
}

func readFileMetadata(metadata []Block) string {
	response := ""
	for index, block := range metadata {
		if index > 0 {
			response = response + " "
		}
		response = response + block.Adress
	}
	return response
}
