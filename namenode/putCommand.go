package main

import (
	"dfs/utils"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sort"
	"strconv"
)

func putCommand(data utils.Message) {
	if fileOnDFS(data.Parameters[0]) {
		logger.Printf("INFO: Se hizo un put sobre el archivo \"%s\" que esta en el dfs, solicitando confirmacion de sobreescritura.\n", data.Parameters[0])
		if !userWantsOverWriteFile(data.Connection) {
			logger.Printf("INFO: El usuario no quiere sobreescribir el archivo \"%s\", abortando comando.\n", data.Parameters[0])
			return
		}
		logger.Printf("INFO: El usuario confirmo la sobreescritura sobre el archivo \"%s\".\n", data.Parameters[0])
		//removeFileData(data.Parameters[0])
	}
	cantBlocks := getCantBlocks(data.Parameters[1])
	if cantBlocks == -1 {
		logger.Print("ERROR: No se pudo realizar la conversion de string a entero.\n")
		return
	}
	assignBlocks(data)
}

func fileOnDFS(fileName string) bool {
	metadataMutex.RLock()
	fileOnDFS := metadata[fileName] != nil
	metadataMutex.RUnlock()
	return fileOnDFS
}

func userWantsOverWriteFile(conn net.Conn) bool {
	overWriteMessage := utils.Message{
		Connection: conn,
		Command:    "overwrite",
	}
	overWriteMessage.Send()
	overWriteResponse := utils.ReadMessage(conn, logger)
	return overWriteResponse.Command == "confirm"
}

func getCantBlocks(blocks string) int {
	cantBlocks, err := strconv.Atoi(blocks)
	if err != nil {
		return -1
	}
	return cantBlocks
}

func assignBlocks(data utils.Message) {
	cantBlocks := getCantBlocks(data.Parameters[1])
	fmt.Print("Solicitando el lock de la info del sistema\n")
	systemInfoMutex.Lock()
	fmt.Print("Lock de la info del sistema obtenido\n")
	orderedBlocks := make([]BlockInfo, len(systemInfo.BlockUsage))
	copy(orderedBlocks, systemInfo.BlockUsage)

	sort.Slice(orderedBlocks, func(i, j int) bool {
		return orderedBlocks[i].Usage < orderedBlocks[j].Usage
	})

	response := utils.Message{
		Connection: data.Connection,
		Command:    "addresses",
		Parameters: make([]string, cantBlocks),
	}

	blocks := make([]Block, cantBlocks)
	for i := 1; i <= cantBlocks; i++ {
		response.Parameters[i-1] = orderedBlocks[0].Address
		blocks[i-1] = Block{
			Name:   fmt.Sprintf("b%d", i),
			Adress: orderedBlocks[0].Address,
		}
		incrementFirstElement(orderedBlocks)
	}
	response.Send()

	blocksSaved := getPreAssignResponse(data.Connection)
	if blocksSaved {
		saveBlockUsageModifications(orderedBlocks)
		systemInfoMutex.Unlock()
		fmt.Print("Lock de la info del sistema liberado\n")
		fmt.Print("Solicitando el lock de la metadata\n")
		metadataMutex.Lock()
		fmt.Print("Lock de la metadata obtenido\n")
		addFileToMetadata(data.Parameters[0], blocks)
		metadataMutex.Unlock()
		fmt.Print("Lock de la metadata liberado\n")
	} else {
		systemInfoMutex.RUnlock()
		fmt.Print("Lock de la info del sistema liberado\n")
	}
}

func incrementFirstElement(orderedBlocks []BlockInfo) {
	orderedBlocks[0].Usage++
	moveIndex := 0
	for moveIndex < len(orderedBlocks)-1 && orderedBlocks[moveIndex].Usage > orderedBlocks[moveIndex+1].Usage {
		aux := orderedBlocks[moveIndex]
		orderedBlocks[moveIndex] = orderedBlocks[moveIndex+1]
		orderedBlocks[moveIndex+1] = aux
		moveIndex++
	}
}

func getPreAssignResponse(conn net.Conn) bool {
	response := utils.ReadMessage(conn, logger)
	return response.Command == "confirm"
}

func saveBlockUsageModifications(newBlocks []BlockInfo) {
	systemInfo.BlockUsage = newBlocks
	file, errCreate := os.Create("data/system_info.json")
	if errCreate != nil {
		panic(errCreate)
	}
	defer file.Close()
	jsonData, errMarshal := json.MarshalIndent(systemInfo, "", "\t")
	if errMarshal != nil {
		panic(errMarshal)
	}

	file.Write(jsonData)
}

func addFileToMetadata(fileName string, blocks []Block) {
	metadata[fileName] = blocks
	file, _ := os.Create("data/metadata.json")
	jsonData, _ := json.MarshalIndent(metadata, "", "\t")
	file.Write(jsonData)
}
