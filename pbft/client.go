package pbft

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"
)

var SendCount int
var GetCount int

func ClientSendMessageAndListen(nodeTable map[string]string) {
	ClientMsgMap = make(map[int]ClientMsg)
	NodeTable = nodeTable
	timestamp1 := time.Now().Unix()
	//开启客户端的本地监听（主要用来接收节点的reply信息）
	go ClientTcpListen()
	fmt.Printf("客户端开启监听，地址：%s\n", clientAddr)

	//stdReader := bufio.NewReader(os.Stdin)
	for {
		data := "hello world!"
		// data, err := stdReader.ReadString('\n')
		// if err != nil {
		// 	fmt.Println("Error reading from stdin")
		// 	panic(err)
		// }
		r := new(Request)
		r.Timestamp = time.Now().UnixNano()
		r.ClientAddr = clientAddr
		r.Message.ID = getRandom()
		//消息内容就是用户的输入
		r.Message.Content = strings.TrimSpace(data)
		br, err := json.Marshal(r)
		if err != nil {
			log.Panic(err)
		}
		fmt.Println(string(br))
		content := jointMessage(cRequest, br)
		//默认N0为主节点，直接把请求信息发送至N0
		tcpDial(content, NodeTable["N0"])
		SendCount++

		if time.Now().Unix()-timestamp1 > 60*1000 {
			break
		}
	}

	fmt.Printf("发送连接数：%d\n", SendCount)
	fmt.Printf("接受连接数：%d\n", GetCount)
}

//返回一个十位数的随机数，作为msgid
func getRandom() int {
	x := big.NewInt(10000000000)
	for {
		result, err := rand.Int(rand.Reader, x)
		if err != nil {
			log.Panic(err)
		}
		if result.Int64() > 1000000000 {
			return int(result.Int64())
		}
	}
}

func handleReply(data []byte) {
	cmd, content := splitMessage(data)
	if cmd != "reply" {
		return
	}
	r := new(Reply)
	err := json.Unmarshal(content, r)
	if err != nil {
		log.Panic(err)
	}

	clientM, ok := ClientMsgMap[r.MessageID]
	if ok {
		clientM.Count++
		if clientM.Count == len(NodeTable) {
			delete(ClientMsgMap, r.MessageID)
			GetCount++
		}
	} else {
		clientM := new(ClientMsg)
		clientM.NodeID = r.NodeID
		clientM.Count = 1
		clientM.Content = r.Content
		ClientMsgMap[r.MessageID] = *clientM
	}

}
