package main

import (
	"dfs/utils"
	"fmt"
	"net"
	"os"
)

func putCommand(data utils.Message) {
	fileBlocks := divideFile(data)
	namenodeConn, groupedAddresses := getFilePartsAddress(data, fileBlocks)
	sendBlocksToDatanodes(data.Parameters[0], fileBlocks, groupedAddresses)
	utils.SendMessage(namenodeConn, "confirm")
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
	message := fmt.Sprintf("put %s %d", data.Parameters[0], len(fileBlocks)+1)
	utils.SendMessage(conn, message)
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
		message := fmt.Sprintf("store %s %d", fileName, len(blocksIndexArray))
		utils.SendMessage(conn, message)
		sendBlocks(conn, blocks, blocksIndexArray)
	}
}

func sendBlocks(conn net.Conn, blocks []FilePart, blocksIndexArray []int) {
	for _, blockIndex := range blocksIndexArray {
		response := utils.ReadMessage(conn, nil)
		if response.Command == "block" {
			message := fmt.Sprintf("blockNum %d", blockIndex)
			utils.SendMessage(conn, message)
			response = utils.ReadMessage(conn, nil)
			if response.Command == "sendData" {
				utils.SendData(conn, blocks[blockIndex-1].Data)
			}
		}
	}
}
