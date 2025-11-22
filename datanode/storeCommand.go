package main

import (
	"dfs/utils"
	"fmt"
	"os"
	"strconv"
)

func storeCommand(data utils.Message) {
	cantBlocks := getCantBlocks(data.Parameters[1])
	buffer := make([]byte, 1024)
	for range cantBlocks {
		message := utils.Message{
			Connection: data.Connection,
			Command:    "block",
		}
		utils.SendMessage(message)
		response := utils.ReadMessage(data.Connection, logger)
		if response.Command == "blockNum" {
			block := createBlock(data.Parameters[0], response.Parameters[0])
			message.Command = "sendData"
			utils.SendMessage(message)
			readedBytes := utils.ReadData(data.Connection, buffer, logger)
			block.Write(buffer[:readedBytes])
		}
	}
}

func getCantBlocks(num string) int {
	n, err := strconv.Atoi(num)
	if err != nil {
		n = -1
	}
	return n
}

func createBlock(fileName string, blockNum string) *os.File {
	fullFileName := fmt.Sprintf("%s-block_%s.%s", fileName[:len(fileName)-4], blockNum, fileName[len(fileName)-3:])
	fullPath := fmt.Sprintf("blocks/%s", fullFileName)
	file, _ := os.Create(fullPath)
	return file
}
