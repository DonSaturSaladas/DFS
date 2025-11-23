package main

import (
	"dfs/utils"
	"fmt"
	"net"
	"os"
	"strconv"
)

func putCommand(data utils.Message) {
	if len(data.Parameters) < 1 || data.Parameters[0] == "" {
		fmt.Print("ERROR: Ingrese el nombre del archivo\n")
		logger.Print("ERROR: Se ingreso el comando put sin el nombre del archivo.\n")
		return
	}
	fileName := data.Parameters[0]
	file, err := openFile(fileName)
	defer file.Close()
	logger.Printf("INFO: Abriendo el archivo \"%s\".\n", fileName)
	if err {
		return
	}
	logger.Printf("INFO: Archivo \"%s\" abierto correctamente.\n", fileName)
	logger.Printf("INFO: Dividiendo el archivo \"%s\".\n", fileName)
	fileBlocks := divideFile(file)
	logger.Printf("INFO: El archivo \"%s\" se dividio en %d bloques.\n", fileName, len(fileBlocks))
	logger.Printf("INFO: Solicitando las direcciones para guardar los bloques del archivo \"%s\".", fileName)
	namenodeConn, groupedAddresses := getFilePartsAddress(data, fileBlocks)
	logger.Printf("INFO: Obtenidas las direcciones para guardar los bloques del archivo \"%s\".", fileName)
	logger.Printf("INFO: Enviando los bloques del archivo \"%s\" a los datanodes correspondientes.\n", fileName)
	sendBlocksToDatanodes(fileName, fileBlocks, groupedAddresses)
	logger.Printf("INFO: Bloques del archivo \"%s\" enviados.\n", fileName)
	confirmMessage := utils.Message{
		Connection: namenodeConn,
		Command:    "confirm",
	}
	logger.Printf("INFO: Confirmando el envio de los bloques del archivo \"%s\" al namenode.\n", fileName)
	confirmMessage.Send()
	namenodeConn.Close()
}

func openFile(filePath string) (*os.File, bool) {
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("ERROR: El archivo %s no existe.\n", filePath)
		} else if os.IsPermission(err) {
			fmt.Printf("ERROR: Permisos insuficientes para abrir el archivo %s.\n", filePath)
		} else {
			fmt.Printf("ERROR: No se pudo abrir el archivo %s.\n", filePath)
		}
		return nil, true
	}
	return file, false
}

func divideFile(file *os.File) []FilePart {
	fileBlocks := make([]FilePart, 0)
	buffer := make([]byte, 1024)
	var readedBytes = 1024
	blockNumber := 1
	for readedBytes == 1024 {
		readedBytes, _ = file.Read(buffer)
		newPart := FilePart{
			BlockNum: blockNumber,
			Data:     append([]byte(nil), buffer[:readedBytes]...),
		}
		fileBlocks = append(fileBlocks, newPart)
		blockNumber++
	}
	return fileBlocks
}

func getFilePartsAddress(data utils.Message, fileBlocks []FilePart) (net.Conn, map[string][]int) {
	conn := connectToNode(namenode_ip)
	message := utils.Message{
		Connection: conn,
		Command:    "put",
		Parameters: []string{data.Parameters[0], strconv.Itoa(len(fileBlocks))},
	}
	message.Send()
	response := utils.ReadMessage(conn, nil)
	if response.Command == "addresses" {
		assignAddressesToBlocks(fileBlocks, response.Parameters)
	}
	groupedAddresses := groupAddresses(data.Parameters[0], response.Parameters)
	return conn, groupedAddresses
}

func assignAddressesToBlocks(blocks []FilePart, addresses []string) {
	for index, address := range addresses {
		blocks[index].DestNode = address
	}
}

func sendBlocksToDatanodes(fileName string, blocks []FilePart, groupedAddresses map[string][]int) {
	for address, blocksIndexArray := range groupedAddresses {
		conn := connectToNode(address)
		message := utils.Message{
			Connection: conn,
			Command:    "store",
			Parameters: []string{fileName, strconv.Itoa(len(blocksIndexArray))},
		}

		message.Send()
		logger.Printf("INFO: Enviando los bloques del archivo \"%s\" al nodo %s.\n", fileName, address)
		sendBlocks(conn, fileName, blocks, blocksIndexArray)
		conn.Close()
	}
}

func sendBlocks(conn net.Conn, fileName string, blocks []FilePart, blocksIndexArray []int) {
	for _, blockIndex := range blocksIndexArray {
		response := utils.ReadMessage(conn, nil)
		if response.Command == "block" {
			message := utils.Message{
				Connection: conn,
				Command:    "blockNum",
				Parameters: []string{strconv.Itoa(blockIndex)},
			}
			message.Send()
			response = utils.ReadMessage(conn, nil)
			if response.Command == "sendData" {
				logger.Printf("INFO: Enviando el bloque %d del archivo \"%s\".\n", blockIndex, fileName)
				utils.SendData(conn, blocks[blockIndex-1].Data)
			}
		}
	}
}
