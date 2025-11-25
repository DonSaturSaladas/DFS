package main

import (
	"dfs/utils"
	"fmt"
	"os"
	"strconv"
)

func storeCommand(data utils.Message) {
	cantBlocks := getCantBlocks(data.Parameters[1])
	logger.Printf("INFO: El cliente desea guardar %d bloques en este nodo.\n", cantBlocks)
	buffer := make([]byte, 1024)
	for range cantBlocks {
		message := utils.Message{
			Connection: data.Connection,
			Command:    "block",
		}
		logger.Printf("INFO: Solicitando al cliente el numero de bloque del archivo \"%s\" que desea guardar.\n", data.Parameters[0])
		message.Send()
		response := utils.ReadMessage(data.Connection, logger)
		if response.Command == "blockNum" {
			logger.Printf("INFO: El cliente desea guardar el bloque %s del archivo \"%s\" .\n", response.Parameters[0], data.Parameters[0])
			logger.Printf("INFO: Creando el bloque %s del archivo del archivo \"%s\" .\n", response.Parameters[0], data.Parameters[0])
			block := createBlock(data.Parameters[0], response.Parameters[0])
			message.Command = "sendData"
			logger.Printf("INFO: Esperando la informacion del bloque %s del archivo \"%s\" .\n", response.Parameters[0], data.Parameters[0])
			message.Send()
			logger.Printf("INFO: Recibida la informacion del bloque %s del archivo \"%s\" .\n", response.Parameters[0], data.Parameters[0])
			readedBytes := utils.ReadData(data.Connection, buffer, logger)
			logger.Printf("INFO: Guardando la informacion del bloque %s del archivo \"%s\" .\n", response.Parameters[0], data.Parameters[0])
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
