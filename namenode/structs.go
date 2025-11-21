package main

type Block struct {
	Name   string `json:"block"`
	Adress string `json:"node"`
}

type SystemInfo struct {
	NamenodeIp string      `json:"namenode_ip"`
	BlockUsage []BlockInfo `json:"block_usage"`
}

type BlockInfo struct {
	Address string `json:"node"`
	Usage   int    `json:"usage"`
}
