// @Author TrandLiu
// @Date 2024/12/17 16:30:00
// @Desc
package models

import (
	"fmt"
	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"gorm.io/gorm"
	net "net"
	"net/http"
	"strconv"
	"sync"
)

type Message struct {
	gorm.Model
	FromId   uint   //发送者
	TargetId uint   //接收者
	Type     string //消息类型 群聊 私聊 广播
	Media    int    //消息类型 文字图片音频
	Content  string //消息内容
	Pic      string
	Url      string
	Desc     string
	Amount   int //其他数字统计
}

func (table *Message) TableName() string {
	return "message"
}

type Node struct {
	Conn      *websocket.Conn
	DataQueue chan []byte
	GroupSets set.Interface
}

// 映射关系
var clientMap map[int64]*Node = make(map[int64]*Node, 0)

// 读写锁
var rwLocker sync.RWMutex

// 聊天需要 : 发送者ID 接收者ID 消息类型 发送内容 发送类型
func Chat(writer http.ResponseWriter, request *http.Request) {
	//检验token
	query := request.URL.Query()
	Id := query.Get("userId")
	//转换类型
	userId, _ := strconv.ParseInt(Id, 10, 64)
	targetId := query.Get("targetId")
	msgType := query.Get("type")
	context := query.Get("context")
	isvalida := true
	conn, err := (&websocket.Upgrader{
		//token校验
		CheckOrigin: func(r *http.Request) bool {
			return isvalida
		},
	}).Upgrade(writer, request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	//获取Conn
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe),
	}
	//用户关系
	//userid 和Node绑定
	rwLocker.Lock()
	clientMap[userId] = node
	rwLocker.Unlock()
	//5.发送逻辑
	go sendProc(node)
	//6.接受逻辑
	go recvProc(node)

}

func sendProc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}
func recvProc(node *Node) {
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			return
		}
		broadMsg(data)
		fmt.Println("[ws]<<<<", data)
	}
}

var udpsendChan chan []byte = make(chan []byte, 1024)

// 广播消息
func broadMsg(data []byte) {
	udpsendChan <- data
}

func init() {
	go udpSendProc()
	go udpRecbvProc()
}

// udp发送数据协程
func udpSendProc() {
	con, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(192, 168, 121, 134),
		Port: 3000,
	})
	defer con.Close()
	if err != nil {
		fmt.Println(err)
	}
}

// udp接受数据协程
func udpRecbvProc() {

}
