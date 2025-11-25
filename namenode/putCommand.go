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
		logger.Print("INFO: Enviando al usuario las direcciones de los namenodes que contienen los bloques a eliminar del archivo que se desea sobreescribir.\n")
		utils.Message{
			Connection: data.Connection,
			Command:    "remove",
			Parameters: getFileAddresses(data.Parameters[0]),
		}.Send()
		logger.Print("INFO: Esperando confirmacion de la eliminacion de los bloques.\n")
		response := utils.ReadMessage(data.Connection, logger)
		if response.Command == "confirm" {
			logger.Print("INFO: Se recibio la confirmacion de la eliminacion de los bloques, eliminando los metadatos asociados al archivo viejo.\n")
			removeFileData(data.Parameters[0])
		}
	}
	cantBlocks := getCantBlocks(data.Parameters[1])
	if cantBlocks == -1 {
		logger.Print("ERROR: No se pudo realizar la conversion de string a entero.\n")
		return
	}
	assignBlocks(data)
}

func fileOnDFS(fileName string) bool {
	logger.Print("INFO: Solicitando mutex para leer la informacion de \"metadata.json\".\n")
	metadataMutex.RLock()
	logger.Print("INFO: Obtenido mutex de \"metadata.json\".\n")
	fileOnDFS := metadata[fileName] != nil
	metadataMutex.RUnlock()
	logger.Print("INFO: Liberado mutex de \"metadata.json\".\n")
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
	logger.Print("INFO: Solicitando mutex para escribir el archivo \"system_info.json\".\n")
	systemInfoMutex.Lock()
	logger.Print("INFO: Obtenido mutex de \"system_info.json\".\n")

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

	logger.Printf("INFO: Asignando nodos a los bloques del archivo \"%s\".\n", data.Parameters[0])
	blocks := make([]Block, cantBlocks)
	for i := 1; i <= cantBlocks; i++ {
		response.Parameters[i-1] = orderedBlocks[0].Address
		blocks[i-1] = Block{
			Name:   fmt.Sprintf("b%d", i),
			Adress: orderedBlocks[0].Address,
		}
		incrementFirstElement(orderedBlocks)
	}
	logger.Printf("INFO: Enviando direcciones de los bloques del archivo \"%s\".\n", data.Parameters[0])
	response.Send()

	logger.Printf("INFO: Esperando confirmacion de que los bloques del archivo \"%s\" fueron guardados en los nodos.\n", data.Parameters[0])
	blocksSaved := getPreAssignResponse(data.Connection)
	if blocksSaved {
		logger.Printf("INFO: Recibida confirmacion de que los bloques del archivo \"%s\" fueron guardados en los nodos.\n", data.Parameters[0])
		logger.Printf("INFO: Guardando las modificaciones asociadas al archivo \"%s\" en el sistema.\n", data.Parameters[0])
		saveBlockUsageModifications(orderedBlocks)

		systemInfoMutex.Unlock()
		logger.Print("INFO: Liberado mutex de \"system_info.json\".\n")
		logger.Print("INFO: Solicitando mutex para escribir en el archivo \"metadata.json\".\n")
		metadataMutex.Lock()
		logger.Print("INFO: Obtenido mutex de \"metadata.json\".\n")
		logger.Printf("INFO: Agregando el archivo \"%s\" a la metadata del sistema.\n", data.Parameters[0])
		addFileToMetadata(data.Parameters[0], blocks)
		metadataMutex.Unlock()
		logger.Print("INFO: Liberado mutex de \"metadata.json\".\n")
	} else {
		logger.Printf("INFO: No se guardaron los bloques del archivo \"%s\" descartando modificaciones realizadas.\n", data.Parameters[0])
		systemInfoMutex.Unlock()
		logger.Print("INFO: Liberado mutex de \"system_info.json\".\n")
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
