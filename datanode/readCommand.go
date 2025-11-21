package main

import (
	"dfs/utils"
	"fmt"
	"os"
)

func readCommand(data utils.Message) {
	cantBlocks := getCantBlocks(data.Parameters[1])
	buffer := make([]byte, 1024)
	for range cantBlocks {
		utils.SendMessage(data.Connection, "block")
		response := utils.ReadMessage(data.Connection, logger)
		if response.Command == "blockNum" {
			fileName := getFullFileName(data.Parameters[0], response.Parameters[0])
			readedBytes := readBlock(fileName, buffer)
			fmt.Printf("Contenido del archivo: \n%s\n", string(buffer))
			utils.SendData(data.Connection, buffer[:readedBytes])
		}
	}
}

func getFullFileName(fileName string, blockNum string) string {
	return fmt.Sprintf("blocks/%s-block_%s.%s", fileName[:len(fileName)-4], blockNum, fileName[len(fileName)-3:])
}

func readBlock(fileName string, buffer []byte) int {
	file, _ := os.Open(fileName)
	readedBytes, _ := file.Read(buffer)
	return readedBytes
}
