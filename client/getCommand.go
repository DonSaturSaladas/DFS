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
		fmt.Print("ERROR: Ingrese el nombre del archivo")
	}
	addresses := getFileAdresses(data.Parameters[0])
	groupedAddresses := groupAddresses(addresses.Parameters)
	fileParts := make([]FilePart, len(addresses.Parameters))
	retrieveFileData(data, fileParts, groupedAddresses)
	fmt.Print(len(fileParts))
	for _, part := range fileParts {
		fmt.Printf("Parte %d: %s", part.BlockNum, string(part.Data))
	}
}

func getFileAdresses(fileName string) utils.Message {
	conn := connectToNode(namenode_ip)
	return getFileBlocks(conn, fileName)
}

func connectToNode(ip string) net.Conn {
	conn, _ := net.Dial("tcp", ip)
	return conn
}

func getFileBlocks(conn net.Conn, fileName string) utils.Message {
	message := utils.Message{
		Connection: conn,
		Command:    "get",
		Parameters: []string{fileName},
	}
	utils.SendMessage(message)
	response := utils.ReadMessage(conn, nil)
	return response
}

func groupAddresses(addresses []string) map[string][]int {
	addressMap := map[string][]int{}
	for blockNum, address := range addresses {
		addressMap[address] = append(addressMap[address], blockNum+1)
	}
	return addressMap
}

func retrieveFileData(data utils.Message, fileParts []FilePart, addresses map[string][]int) {
	var datanodeConn net.Conn
	fileName := data.Parameters[0]
	var message utils.Message
	var response utils.Message
	filePartsIndex := 0

	for ip, blockArray := range addresses {
		datanodeConn = connectToNode(ip)
		cantBlocks := len(addresses[ip])
		message = utils.Message{
			Connection: datanodeConn,
			Command:    "read",
			Parameters: []string{fileName, strconv.Itoa(cantBlocks)},
		}
		utils.SendMessage(message)
		for _, blockIndex := range blockArray {
			buffer := make([]byte, 1024)
			response = utils.ReadMessage(datanodeConn, nil)
			if response.Command == "block" {
				message.Command = "blockNum"
				message.Parameters = []string{strconv.Itoa(blockIndex)}

				utils.SendMessage(message)
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
