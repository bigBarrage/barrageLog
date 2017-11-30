package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"gopkg.in/mgo.v2/bson"
)

var Lock sync.Mutex

func main() {
	Process()
}

func Process() {
	conn, err := net.Dial("tcp", broadcastingStationHost+":"+broadcastingStationPort)
	if err != nil {
		fmt.Println("[连接失败]：", err)
		return
	}
	u := url.URL{}
	u.Host = broadcastingStationHost + ":" + broadcastingStationPort
	u.Path = broadcastingStationUri
	u.Scheme = "ws"

	h := http.Header{}
	c, _, err := websocket.NewClient(conn, &u, h, 1024, 1024)
	if err != nil {
		fmt.Println("[websocket连接失败]：", err)
		return
	}
	for {
		_, reader, err := c.NextReader()
		if err != nil {
			fmt.Println(err)
			return
		}
		msg := make([]byte, 1024)
		l, _ := reader.Read(msg)
		Lock.Lock()
		fmt.Println("=======================")
		fmt.Println("收到：", string(msg[:l]))

		fmt.Println(1)
		m := make(map[string]interface{})
		json.Unmarshal(msg[:l], &m)

		fmt.Println(2)
		conn, err := getConn()
		fmt.Println(3)
		if err != nil {
			fmt.Println("错误3：", err)
			continue
		}
		fmt.Println(4)
		err = conn.DB(mongoDatabase).C("chatLog_" + time.Now().Format("2006-01-02")).Insert(bson.M(m))
		fmt.Println(5)
		fmt.Println(err)
		fmt.Println(6)
		conn.Close()

		Lock.Unlock()
	}
}

type Message struct {
	MessageType int64  `json:"type"`
	RoomID      string `json:"room_id"`
	Time        int64  `json:"time"`
	Body        string `json:"body"`
}
