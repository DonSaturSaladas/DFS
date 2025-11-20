package main

import (
	"dfs/utils"
)

func lsCommand(params []string) {
	var metadata map[string][]Block
	parseJson(&metadata)

	response := readMetadataFiles(metadata)

	utils.SendMessage(conn, response)
}

func readMetadataFiles(metadata map[string][]Block) string {
	response := "["
	for fileName := range metadata {
		response = response + fileName + ","
	}
	response = response[:len(response)-1] + "]"
	return response
}
