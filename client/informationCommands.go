package main

import (
	"dfs/utils"
	"fmt"
)

func lsCommand(data utils.Message) {
	conn := connectToNode(namenode_ip)
	defer conn.Close()
	message := utils.Message{
		Connection: conn,
		Command:    "ls",
	}
	message.Send()
	response := utils.ReadMessage(conn, nil)
	if response.Command == "files" {
		fmt.Printf("Archivos en el DFS:\n -")
		for i, file := range response.Parameters {
			fmt.Print(" " + file)
			if i != len(response.Parameters)-1 {
				fmt.Print(",")
			}
		}
		fmt.Print(".\n")
	}
}
func infoCommand(data utils.Message) {
	if len(data.Parameters) < 1 || data.Parameters[0] == "" {
		fmt.Print("ERROR: Ingrese el nombre del archivo\n")
		return
	}
	fileInfo := getFileInfo(data)
	showFileInfo(fileInfo)
}

func getFileInfo(data utils.Message) []string {
	conn := connectToNode(namenode_ip)
	defer conn.Close()
	message := utils.Message{
		Connection: conn,
		Command:    "get",
		Parameters: data.Parameters,
	}
	message.Send()
	response := utils.ReadMessage(conn, nil)
	return response.Parameters
}

func showFileInfo(fileInfo []string) {
	if len(fileInfo) == 0 {
		fmt.Print("El archivo no se encuentra en el DFS.\n")
	} else {
		fmt.Print("Las partes del archivo esta en los siguientes nodos:\n")
	}
	for i, address := range fileInfo {
		fmt.Printf("\t- Bloque %d: %s\n", i+1, address)
	}
}
