package main

import (
	"dfs/utils"
	"fmt"
	"net"
	"os"
	"sort"
	"strconv"
)

type FilePart struct {
	Data     []byte
	BlockNum int
	DestNode string
}

func getCommand(data utils.Message) {
	if len(data.Parameters) != 1 || data.Parameters[0] == "" {
		fmt.Print("ERROR: Ingrese el nombre del archivo\n")
		logger.Printf("ERROR: Se ingreso el comando get sin el nombre del archivo.\n")
		return
	}
	logger.Printf("INFO: Solicitando al namenode la direccion de los bloques del archivo %s.\n", data.Parameters[0])
	blockAddresses := getBlocksAdresses(data.Parameters[0])
	defer blockAddresses.Connection.Close()
	if len(blockAddresses.Parameters) == 0 {
		logger.Printf("INFO: El archivo %s no se encontro en el DFS.\n", data.Parameters[0])
		fmt.Printf("El archivo %s no se encuentra en el DFS", data.Parameters[0])
		return
	}
	logger.Printf("INFO: Mensaje con las direcciones de los bloques del archivo %s recibidas.\n", data.Parameters[0])
	groupedAddresses := groupAddresses(data.Parameters[0], blockAddresses.Parameters)
	fileParts := make([]FilePart, len(blockAddresses.Parameters))
	logger.Printf("INFO: Solicitando los bloques del archivo %s.\n", data.Parameters[0])
	retrieveFileData(data, fileParts, groupedAddresses)
	logger.Printf("INFO: Recibidos los bloques del archivo %s.\n", data.Parameters[0])
	saveFile(data.Parameters[0], fileParts)
}

func getBlocksAdresses(fileName string) utils.Message {
	conn := connectToNode(namenode_ip)
	return getFileBlocks(conn, fileName)
}

func connectToNode(ip string) net.Conn {
	logger.Printf("INFO: Conectandose al nodo %s.\n", ip)
	conn, err := net.Dial("tcp", ip)
	if err != nil {
		logger.Printf("ERROR: No se pudo conectar al nodo %s.\n", ip)
		return nil
	}
	logger.Printf("INFO: Conexion establecida con el nodo %s.\n", ip)
	return conn
}

func getFileBlocks(conn net.Conn, fileName string) utils.Message {
	message := utils.Message{
		Connection: conn,
		Command:    "get",
		Parameters: []string{fileName},
	}
	message.Send()
	logger.Printf("INFO: Esperando respuesta del namenode con la direccion de los bloques del archivo %s.\n", fileName)
	response := utils.ReadMessage(conn, nil)
	return response
}

func groupAddresses(fileName string, blockAddresses []string) map[string][]int {
	logger.Printf("INFO: Agrupando los bloques del archivo \"%s\" por nodo.\n", fileName)
	addressMap := map[string][]int{}
	for blockNum, address := range blockAddresses {
		addressMap[address] = append(addressMap[address], blockNum+1)
	}
	logger.Printf("INFO: Los bloques del archivo \"%s\" fueron agrupados correctamente.\n", fileName)
	return addressMap
}

func retrieveFileData(data utils.Message, fileParts []FilePart, blockAddresses map[string][]int) {
	var datanodeConn net.Conn
	fileName := data.Parameters[0]
	var message utils.Message
	var response utils.Message
	filePartsIndex := 0

	for ip, blockArray := range blockAddresses {
		datanodeConn = connectToNode(ip)
		cantBlocks := len(blockAddresses[ip])
		message = utils.Message{
			Connection: datanodeConn,
			Command:    "read",
			Parameters: []string{fileName, strconv.Itoa(cantBlocks)},
		}
		message.Send()
		for _, blockIndex := range blockArray {
			buffer := make([]byte, 1024)
			response = utils.ReadMessage(datanodeConn, nil)
			if response.Command == "block" {
				message.Command = "blockNum"
				message.Parameters = []string{strconv.Itoa(blockIndex)}

				message.Send()
				readedBytes := utils.ReadData(datanodeConn, buffer, nil)
				fileParts[filePartsIndex] = FilePart{
					Data:     buffer[:readedBytes],
					BlockNum: blockIndex,
				}
				filePartsIndex++
			}
		}
	}
}

func saveFile(fileName string, fileParts []FilePart) {
	sortFileParts(fileParts)
	fileData := joinFileParts(fileParts)
	writeDataToFile(fileName, fileData)
}

func sortFileParts(fileParts []FilePart) {
	sort.Slice(fileParts, func(i, j int) bool {
		return fileParts[i].BlockNum < fileParts[j].BlockNum
	})
}

func joinFileParts(fileParts []FilePart) []byte {
	var fileData []byte
	for _, filePart := range fileParts {
		fileData = append(fileData, filePart.Data...)
	}
	return fileData
}

func writeDataToFile(fileName string, fileData []byte) {
	filePath := "downloads/" + fileName
	file, err := os.Create(filePath)
	if err != nil {
		logger.Printf("ERROR: No se pudo crear el archivo \"%s\".\n", fileName)
	} else {
		logger.Printf("INFO: Archivo \"%s\" creado correctamente, escribiendo el contenido.\n", fileName)
		_, err := file.Write(fileData)
		if err != nil {
			logger.Printf("ERROR: No se pudo escribir el contenido del archivo \"%s\", eliminando el archivo.\n", fileName)
			os.Remove(filePath)
		} else {
			logger.Printf("INFO: El archivo \"%s\" fue guardado exitosamente.\n", fileName)
		}
	}
}
