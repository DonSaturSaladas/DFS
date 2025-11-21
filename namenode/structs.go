package main

import "net"

type MessageData struct {
	Connection net.Conn
	Parameters []string
}

type Block struct {
	Name   string `json:"block"`
	Adress string `json:"node"`
}

type SystemInfo struct {
	BlockUsage []BlockInfo `json:"block_usage"`
}

type BlockInfo struct {
	Address string `json:"node"`
	Usage   int    `json:"usage"`
}
