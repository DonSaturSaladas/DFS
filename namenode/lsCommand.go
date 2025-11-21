package main

import (
	"dfs/utils"
)

func lsCommand(data MessageData) {
	var metadata map[string][]Block
	parseJson(&metadata)

	response := readMetadataFiles(metadata)

	utils.SendMessage(data.Connection, response)
}

func readMetadataFiles(metadata map[string][]Block) string {
	response := "["
	for fileName := range metadata {
		response = response + fileName + ","
	}
	return response[:len(response)-1] + "]"
}
