package main


import (
	"flag"
	"net/http"
	"net"
	"github.com/gorilla/websocket"
	"log"
	"fmt"
	"sync"
)

var addr = flag.String("addr", "127.0.0.1:8080", "http service address")
var upgrader = websocket.Upgrader{
	// 解决跨域问题
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize: 65535,		
	WriteBufferSize: 65535,
}

func accept(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("upgrade error:", err)
		return
	} 
	ConnMgrInstance().OnAccept(c)
}

func main() {
	var wg sync.WaitGroup

	MsgQueueMgrInstance().StartMsgQueue()
	ConnMgrInstance()

	conn, err := net.Dial("tcp", "127.0.0.1:10086")
	if nil != err {
		fmt.Println(err.Error())
	} else {
		fmt.Println("conn ok!!!")
		conn.Close()
	}

	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/", accept)

	http.ListenAndServe(*addr, nil)
	wg.Add(1)

	go func() {
		defer wg.Done()
		err1 := http.ListenAndServe(*addr, nil)
		fmt.Printf("Start ListenAndServe, error[%s]\n", err1.Error())
	}()
	
	wg.Wait()
}
