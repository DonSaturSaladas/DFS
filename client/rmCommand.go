package main

import (
	"dfs/utils"
	"fmt"
	"net"
)

func rmCommand(data utils.Message) {
	if len(data.Parameters) < 1 || data.Parameters[0] == "" {
		fmt.Print("ERROR: Ingrese el nombre del archivo\n")
		logger.Print("ERROR: Se ingreso el comando put sin el nombre del archivo.\n")
		return
	}
	logger.Printf("INFO: Solicitando las direcciones de los bloques del archivo \"%s\" para eliminarlas.\n", data.Parameters[0])
	response := getAddressesToRemove(data)
	defer response.Connection.Close()
	if response.Command == "addresses" {
		logger.Printf("INFO: Se recibieron las direcciones para eliminar el archivo \"%s\".\n", data.Parameters[0])
		removeBlocks(data, response.Parameters)
	} else if response.Command == "nofile" {
		logger.Printf("INFO: El archivo \"%s\" no se encuentra en el DFS", data.Parameters[0])
		fmt.Printf("El archivo \"%s\" no se encuentra en el DFS.\n", data.Parameters[0])
	}
}

func getAddressesToRemove(data utils.Message) utils.Message {
	namenodeConn := connectToNode(namenode_ip)
	message := utils.Message{
		Connection: namenodeConn,
		Command:    "rm",
		Parameters: data.Parameters,
	}
	message.Send()
	return utils.ReadMessage(data.Connection, logger)
}

func removeBlocks(data utils.Message, addresses []string) {
	logger.Printf("INFO: Eliminando los bloques del archivo \"%s\".\n", data.Parameters[0])
	removeNodeBlocks(data.Parameters[0], addresses)
	confirmDeletedBlocks(data.Connection)
	logger.Printf("INFO: Se eliminaron los bloques del archivo \"%s\".\n", data.Parameters[0])
}

func removeNodeBlocks(fileName string, blocksAddresses []string) {
	for _, address := range blocksAddresses {
		datanodeConn := connectToNode(address)
		message := utils.Message{
			Connection: datanodeConn,
			Command:    "rm",
			Parameters: []string{fileName},
		}
		logger.Printf("INFO: Enviando mensaje para eliminar el archivo \"%s\" al nodo %s.\n", fileName, address)
		message.Send()
		datanodeConn.Close()
	}
}

func confirmDeletedBlocks(namenodeConn net.Conn) {
	message := utils.Message{
		Connection: namenodeConn,
		Command:    "confirm",
	}
	message.Send()
}
