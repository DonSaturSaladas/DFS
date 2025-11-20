package main

import (
	"dfs/utils"
	"encoding/json"
	"fmt"
	"os"
)

type Block struct {
	Name   string `json:"block"`
	Adress string `json:"node"`
}

func getCommand(params []string) {
	if checkEmptyParams(params, "get") {
		return
	}
	fileName := params[0]
	var metadata map[string][]Block
	parseJson(&metadata)

	fileMetadata := metadata[fileName]
	response := readFileMetadata(fileMetadata)
	utils.SendMessage(conn, response)
}

func checkEmptyParams(params []string, command string) bool {
	if len(params) == 0 {
		logger.Printf("ERROR: El cliente ingreso el comando %s sin parametros.\n", command)
		utils.SendMessage(conn, fmt.Sprintf("Ingrese los parametros luego del comando \"%s\"\n", command))
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
	response := "["
	for index, block := range metadata {
		if index > 0 {
			response = response + ","
		}
		response = response + block.Adress
	}
	response = response + "]"
	return response
}
