package main

import (
	"dfs/utils"
	"fmt"
	"strconv"
)

func putCommand(params []string) {
	if checkPutErrors(params) {
		return
	}
	cantBlocks := getCantBlocks(params[1])
	if cantBlocks == -1 {
		return
	}
	fileName := params[0]

	fmt.Printf("El cliente quiere guardar el archivo %s con %d bloques", fileName, cantBlocks)
}

func checkPutErrors(params []string) bool {
	if checkEmptyParams(params, "put") {
		return true
	}
	if len(params) < 2 {
		logger.Print("ERROR: El cliente ingreso el comando put con menos de 2 parametros.\n")
		utils.SendMessage(conn, "ERROR: El uso del comando es put <nombreArchivo> <cantidadBloques>.\n")
		return true
	}
	if params[0] == "" {
		logger.Print("ERROR: El cliente ingreso el comando put con el nombre de archivo invalido.\n")
		utils.SendMessage(conn, "ERROR: Ingrese un nombre de archivo valido.\n")
		return true
	}
	return false
}

func getCantBlocks(blocks string) int {
	cantBlocks, err := strconv.Atoi(blocks)
	if err != nil {
		logger.Print("ERROR: No se pudo realizar la conversion de string a entero.\n")
		utils.SendMessage(conn, "ERROR: Ingrese un numero valido de bloques.\n")
		return -1
	}
	return cantBlocks
}
