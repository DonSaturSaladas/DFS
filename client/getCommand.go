package main

import (
	"dfs/utils"
	"fmt"
	"net"
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
	groupedAddresses := groupAddresses(blockAddresses.Parameters)
	fileParts := make([]FilePart, len(blockAddresses.Parameters))
	logger.Printf("INFO: Solicitando los bloques del archivo %s.\n", data.Parameters[0])
	retrieveFileData(data, fileParts, groupedAddresses)
	logger.Printf("INFO: Recibidos los bloques del archivo %s.\n", data.Parameters[0])
	fmt.Print(len(fileParts))
	for _, part := range fileParts {
		fmt.Printf("Parte %d: %s", part.BlockNum, string(part.Data))
	}
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

func groupAddresses(blockAddresses []string) map[string][]int {
	addressMap := map[string][]int{}
	for blockNum, address := range blockAddresses {
		addressMap[address] = append(addressMap[address], blockNum+1)
	}
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
