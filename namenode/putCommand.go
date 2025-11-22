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
	if checkPutErrors(data) {
		return
	}
	cantBlocks := getCantBlocks(data.Parameters[1])
	if cantBlocks == -1 {
		logger.Print("ERROR: No se pudo realizar la conversion de string a entero.\n")
		utils.SendMessage(data.Connection, "ERROR: Ingrese un numero valido de bloques.\n")
		return
	}
	preAssignBlocks(data)
}

func checkPutErrors(data utils.Message) bool {
	if checkEmptyParams(data) {
		return true
	}
	if len(data.Parameters) < 2 {
		logger.Print("ERROR: El cliente ingreso el comando put con menos de 2 parametros.\n")
		utils.SendMessage(data.Connection, "ERROR: El uso del comando es put <nombreArchivo> <cantidadBloques>.\n")
		return true
	}
	if data.Parameters[0] == "" {
		logger.Print("ERROR: El cliente ingreso el comando put con el nombre de archivo invalido.\n")
		utils.SendMessage(data.Connection, "ERROR: Ingrese un nombre de archivo valido.\n")
		return true
	}
	return false
}

func getCantBlocks(blocks string) int {
	cantBlocks, err := strconv.Atoi(blocks)
	if err != nil {
		return -1
	}
	return cantBlocks
}

func preAssignBlocks(data utils.Message) {
	cantBlocks, _ := strconv.Atoi(data.Parameters[1])
	orderedBlocks := make([]BlockInfo, len(systemInfo.BlockUsage))
	copy(orderedBlocks, systemInfo.BlockUsage)

	sort.Slice(orderedBlocks, func(i, j int) bool {
		return orderedBlocks[i].Usage < orderedBlocks[j].Usage
	})

	responseString := ""
	blocks := make([]Block, cantBlocks)
	for i := 1; i <= cantBlocks; i++ {
		responseString = responseString + orderedBlocks[0].Address + " "
		blocks[i-1] = Block{
			Name:   fmt.Sprintf("b%d", i),
			Adress: orderedBlocks[0].Address,
		}
		incrementFirstElement(orderedBlocks)
	}
	responseString = responseString[:len(responseString)-1]
	utils.SendMessage(data.Connection, responseString)
	blocksSaved := getPreAssignResponse(data.Connection)
	if blocksSaved {
		saveBlockUsageModifications(orderedBlocks)
		addFileToMetadata(data.Parameters[0], blocks)
	} else {
		//Discard JSON
		fmt.Print("Bloques NO guardados\n")
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
	file, _ := os.Create("data/system_info.json")
	jsonData, _ := json.MarshalIndent(systemInfo, "", "\t")

	file.Write(jsonData)
}

func addFileToMetadata(fileName string, blocks []Block) {
	var metadata map[string][]Block
	parseJson(&metadata)
	metadata[fileName] = blocks
	file, _ := os.Create("data/metadata.json")
	jsonData, _ := json.MarshalIndent(metadata, "", "\t")
	file.Write(jsonData)
}
