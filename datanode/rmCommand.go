package main

import (
	"dfs/utils"
	"os"
	"strings"
)

func rmCommand(data utils.Message) {
	fileName, fileExtension := getFileNameData(data.Parameters[0])
	blocksEntries, _ := os.ReadDir("./blocks")

	for _, entry := range blocksEntries {
		entryName, entryExtension := getFileNameData(entry.Name())
		if strings.Contains(entryName, fileName) && entryExtension == fileExtension {
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
