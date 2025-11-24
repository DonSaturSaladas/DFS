package main

import (
	"bufio"
	"dfs/utils"
	"fmt"
	"os"
)

const namenode_ip = "localhost:5000"

var logger = utils.CreateLogger("CLIENTE")

func main() {
	welcomeUser()
	var userInput utils.Message
	for userInput.Command != "exit" {
		userInput = getUserInput()
		commandsMap := getCommandsMap()
		commandReaded := commandsMap[userInput.Command]
		if commandReaded != nil {
			logger.Printf("INFO: Ejecutando comando \"%s\"\n", userInput.Command)
			commandReaded(userInput)
		} else {
			logger.Printf("INFO: Se ingreso un comando invalido -> \"%s\"\n", userInput.Command)
			welcomeUser()
			fmt.Print("\nIngrese un comando valido")
		}
	}
}

func getUserInput() utils.Message {
	buffer := bufio.NewReader(os.Stdin)
	fmt.Print("\nIngrese un comando: ")
	logger.Print("INFO: Esperando un comando.")
	commandReaded, _ := buffer.ReadString('\n')
	message := utils.ParseMessageData(commandReaded)
	return message
}

func getCommandsMap() map[string]func(utils.Message) {
	return map[string]func(utils.Message){
		"put":  putCommand,
		"get":  getCommand,
		"ls":   lsCommand,
		"info": infoCommand,
		"rm":   rmCommand,
	}
}

func welcomeUser() {
	fmt.Printf("════════════════════      Sistema de archivos distribuido      ════════════════════\n")
	fmt.Printf("║                                                                                 ║\n")
	fmt.Printf("║      Comandos disponibles:                                                      ║\n")
	fmt.Printf("║           - put <archivo> : Ingresa un archivo al DFS                           ║\n")
	fmt.Printf("║                                                                                 ║\n")
	fmt.Printf("║           - get <archivo> : Descarga un archivo del dfs a la carpeta            ║\n")
	fmt.Printf("║                            downloads                                            ║\n")
	fmt.Printf("║                                                                                 ║\n")
	fmt.Printf("║           - ls            : Lista los archivos guardados en el DFS              ║\n")
	fmt.Printf("║                                                                                 ║\n")
	fmt.Printf("║           - info <archivo>: Muestra la informacion de un archivo                ║\n")
	fmt.Printf("║                                                                                 ║\n")
	fmt.Printf("║           - exit          : Salir de la interfaz                                ║\n")
	fmt.Printf("║                                                                                 ║\n")
	fmt.Printf("════════════════════════      Maxi Barco - LU: 142557      ════════════════════════\n")
}
