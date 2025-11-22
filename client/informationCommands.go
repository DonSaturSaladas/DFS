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
	utils.SendMessage(message)
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
	fmt.Print("Comando info")
}
