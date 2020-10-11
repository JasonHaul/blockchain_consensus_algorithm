package main

import (
	"blockchain_consensus_algorithm/pbft"
	"os"
	"strconv"
)

func main() {
	//为四个节点生成公私钥
	pbft.GenRsaKeys()
	NodeTable := make(map[string]string)
	for i := 0; i <= 9; i++ {
		n := "N"
		n += strconv.Itoa(i)
		addr := "127.0.0.1:"
		addr += strconv.Itoa(i + 8000)
		NodeTable[n] = addr
	}
	nodeID := os.Args[1]
	if nodeID == "client" {
		pbft.ClientSendMessageAndListen(NodeTable) //启动客户端程序
	} else if nodeID == "node" {
		pbft.PBFT(NodeTable)
	}
	select {}
}
