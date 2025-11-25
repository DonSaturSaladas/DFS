package main

import (
	"dfs/utils"
	"os"
	"strings"
)

func rmCommand(data utils.Message) {
	fileName, fileExtension := getFileNameData(data.Parameters[0])
	blocksEntries, _ := os.ReadDir("./blocks")
	logger.Printf("INFO: Buscando bloques asociados al archivo \"%s\".\n", data.Parameters[0])
	for _, entry := range blocksEntries {
		entryName, entryExtension := getFileNameData(entry.Name())
		if strings.Contains(entryName, fileName) && entryExtension == fileExtension {
			logger.Printf("INFO: Se encontro un bloque asociado al archivo \"%s\", eliminando bloque.\n", data.Parameters[0])
			removeFile(entry.Name())
		}
	}
}

func getFileNameData(fullName string) (string, string) {
	extensionPointIndex := strings.LastIndex(fullName, ".")
	if extensionPointIndex < 0 {
		return fullName, ""
	}
	return fullName[:extensionPointIndex], fullName[extensionPointIndex:]
}

func removeFile(fileName string) {
	fullPath := "blocks/" + fileName
	os.Remove(fullPath)
}
