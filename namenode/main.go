package main

import (
	"dfs/utils"
	"fmt"
	"log"
	"net"
	"strings"
)

var port int = 5000
var logger *log.Logger

func main() {
	listener := startServer()
	for {
		conn := acceptConnection(listener)
		go handleConnection(conn)
	}
}

func startServer() net.Listener {
	logger = utils.CreateLogger("NAMENODE")
	listener, err := utils.ListenPort(port)
	if err != 0 {
		error_string := fmt.Sprintf("ERROR: No se pudo escuchar en el puerto %d\n", port)
		logger.Print(error_string)
		panic(error_string)
	}
	logger.Printf("INFO: Iniciado el servidor, escuchando el puerto %d\n", port)
	return listener
}

func acceptConnection(listener net.Listener) net.Conn {
	conn, err := listener.Accept()
	if err != nil {
		error_string := "ERROR: No se pudo aceptar la conexion\n"
		logger.Print(error_string)
		panic(error_string)
	}
	return conn
}

func handleConnection(conn net.Conn) {
	type Commands struct {
		Get string
		Put string
		Ls  string
	}
	commands := Commands{
		Get: "get",
		Put: "put",
		Ls:  "ls",
	}
	buffer := make([]byte, 256)
	commandsMap := map[string]func(){
		commands.Get: getCommand,
		commands.Put: putCommand,
		commands.Ls:  lsCommand,
	}

	for {
		readedBytes := readMessage(conn, buffer)
		commandString := strings.TrimSpace(string(buffer[:readedBytes]))
		fmt.Printf("Comando leido: %s\n", commandString)
		commandReaded := commandsMap[commandString]
		if commandReaded == nil {
			fmt.Printf("Comando Null\n")
			conn.Write([]byte("El comando ingresado no es valido.\nIngrese uno de los siguientes comandos: get, put, ls\n"))
		} else {
			commandReaded()
		}
	}
}

func readMessage(conn net.Conn, buffer []byte) int {
	readedBytes, err := conn.Read(buffer)
	fmt.Printf("Ingresado algo\n")
	if err != nil {
		error_string := "ERROR: No se pudo leer en la conexion\n"
		logger.Print(error_string)
		panic(error_string)
	}
	return readedBytes
}

func getCommand() {
	fmt.Print("Get command\n")
}

func putCommand() {
	fmt.Print("Put command\n")
}

func lsCommand() {
	fmt.Print("Ls command\n")
}
