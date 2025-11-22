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
		return
	}
	fileBlocks := divideFile(data)
	namenodeConn, groupedAddresses := getFilePartsAddress(data, fileBlocks)
	sendBlocksToDatanodes(data.Parameters[0], fileBlocks, groupedAddresses)
	confirmMessage := utils.Message{
		Connection: namenodeConn,
		Command:    "confirm",
	}
	utils.SendMessage(confirmMessage)
	namenodeConn.Close()
}

func divideFile(data utils.Message) []FilePart {
	fileName := data.Parameters[0]
	file, _ := os.Open(fileName)

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
	utils.SendMessage(message)
	response := utils.ReadMessage(conn, nil)
	if response.Command == "addresses" {
		assignAddressesToBlocks(fileBlocks, response.Parameters)
	}
	return conn, groupAddresses(response.Parameters)
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

		utils.SendMessage(message)
		sendBlocks(conn, blocks, blocksIndexArray)
		conn.Close()
	}
}

func sendBlocks(conn net.Conn, blocks []FilePart, blocksIndexArray []int) {
	for _, blockIndex := range blocksIndexArray {
		response := utils.ReadMessage(conn, nil)
		if response.Command == "block" {
			message := utils.Message{
				Connection: conn,
				Command:    "blockNum",
				Parameters: []string{strconv.Itoa(blockIndex)},
			}
			utils.SendMessage(message)
			response = utils.ReadMessage(conn, nil)
			if response.Command == "sendData" {
				utils.SendData(conn, blocks[blockIndex-1].Data)
			}
		}
	}
}
