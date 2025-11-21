package main

import (
	"dfs/utils"
	"encoding/json"
	"fmt"
	"os"
)

func getCommand(data utils.Message) {
	if checkEmptyParams(data) {
		return
	}
	fileName := data.Parameters[0]
	var metadata map[string][]Block
	parseJson(&metadata)

	fileMetadata := metadata[fileName]
	response := readFileMetadata(fileMetadata)
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

func parseJson(variable *map[string][]Block) {
	byteValue, err := os.ReadFile("data/metadata.json")
	if err != nil {
		logger.Print("ERROR: No se pudo abrir el archivo de metadata.")
		panic(err)
	}
	json.Unmarshal(byteValue, variable)
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
