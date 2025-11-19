package main

import (
	"dfs/utils"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type Block struct {
	Name   string `json:"block"`
	Adress string `json:"node"`
}

var metadata map[string][]Block

func checkParamsEmpty(conn net.Conn, params string, command string) bool {
	if params == "" {
		logger.Printf("ERROR: El cliente ingreso el comando %s sin parametros.\n", command)
		utils.SendMessage(conn, fmt.Sprintf("Ingrese los parametros luego del comando \"%s\"\n", command))
		return true
	}
	return false
}

func getCommand(conn net.Conn, fileName string) {
	if checkParamsEmpty(conn, fileName, "get") {
		return
	}
	parseJson()

	fileMetadata := metadata[fileName]

	response := "["
	for index, block := range fileMetadata {
		if index > 0 {
			response = response + ","
		}
		response = response + block.Adress
	}
	response = response + "]"

	utils.SendMessage(conn, response)
}

func parseJson() {
	byteValue, err := os.ReadFile("data/metadata.json")
	if err != nil {
		logger.Print("ERROR: No se pudo abrir el archivo de metadata.")
		panic(err)
	}

	json.Unmarshal(byteValue, &metadata)
}

func putCommand(conn net.Conn, params string) {
	if checkParamsEmpty(conn, params, "put") {
		return
	}
}

func lsCommand(conn net.Conn, params string) {
	parseJson()

	response := "["
	writeComma := false
	for fileName := range metadata {
		if !writeComma {
			response = response + fileName
			writeComma = true
		} else {
			response = response + "," + fileName
		}
	}
	response = response + "]"

	utils.SendMessage(conn, response)
}
