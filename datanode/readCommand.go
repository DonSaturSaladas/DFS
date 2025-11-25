package main

import (
	"dfs/utils"
	"fmt"
	"os"
)

func readCommand(data utils.Message) {
	cantBlocks := getCantBlocks(data.Parameters[1])
	logger.Printf("INFO: El cliente desea obtener %d bloques de este nodo.\n", cantBlocks)
	buffer := make([]byte, 1024)
	for range cantBlocks {
		message := utils.Message{
			Connection: data.Connection,
			Command:    "block",
		}
		logger.Printf("INFO: Solicitando el numero de bloque que se desea obtener del archivo \"%s\".\n", data.Parameters[0])
		message.Send()
		response := utils.ReadMessage(data.Connection, logger)
		if response.Command == "blockNum" {
			logger.Printf("INFO: El cliente desea obtener el bloque %s del archivo \"%s\".\n", response.Parameters[0], data.Parameters[0])
			fileName := getFullFileName(data.Parameters[0], response.Parameters[0])
			logger.Printf("INFO: Obteniendo el bloque %s del archivo\"%s\".\n", response.Parameters[0], data.Parameters[0])
			readedBytes := readBlock(fileName, buffer)
			logger.Printf("INFO: Enviando el bloque %s del archivo \"%s\".\n", response.Parameters[0], data.Parameters[0])
			utils.SendData(data.Connection, buffer[:readedBytes])
			logger.Printf("INFO: Bloque %s del archivo \"%s\" enviado.\n", response.Parameters[0], data.Parameters[0])
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
