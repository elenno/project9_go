package main

import (
	"github.com/gorilla/websocket"
	"sync"
	"fmt"
	"io"
	"net"
)

type WSConn struct {
	netid int
	conn *websocket.Conn
}

type TCPConn struct {
	netid int
	conn *net.Conn
}

type ConnMgr struct {
	wsConnMap map[int]WSConn
	tcpConnMap map[int]TCPConn
	gameNetID int
}

var connMgrInst *ConnMgr
var instMutex sync.Mutex

func ConnMgrInstance() *ConnMgr {
	if nil == connMgrInst {
		instMutex.Lock()
		if nil == connMgrInst {
			connMgrInst = &ConnMgr{
				wsConnMap: make(map[int]WSConn),
				tcpConnMap: make(map[int]TCPConn),
				gameNetID: -1,
			}
		}
		instMutex.Unlock()	
	}		
    return connMgrInst
}

func (mgr *ConnMgr) OnAccept(conn *websocket.Conn) {
	netid := GetNewNetID()
	fmt.Printf("Accept WSConn netid[%d] from addr[%s]\n", netid, conn.RemoteAddr().String())

	
	wsConn := WSConn{netid:netid, conn:conn}
	mgr.wsConnMap[netid] = wsConn	
	
	go mgr.StartRead(netid, conn)
}

func (mgr *ConnMgr) OnConnectGameWorld(conn *net.Conn) {
	netid := GetNewNetID()
	fmt.Printf("Connect Gameworld succ, netid[%d] addr[%s]\n", netid, (*conn).RemoteAddr().String())

	tcpConn := TCPConn{netid:netid, conn:conn}
	mgr.tcpConnMap[netid] = tcpConn

	mgr.gameNetID = netid
}

func (mgr *ConnMgr) OnDisconnect(netid int) {
	delete(mgr.wsConnMap, netid)
	fmt.Printf("Disconnect WSConn netid[%d]\n", netid)
}

func (mgr *ConnMgr) StartRead(netid int, conn *websocket.Conn) {
	defer mgr.OnDisconnect(netid)
	defer conn.Close()

	for {
		msgType, buff, err := conn.ReadMessage()
		if nil != err {
			if err != io.EOF {
				fmt.Printf("conn netid[%d] read error[%s]\n", netid, err)
			}
			return
		} else if msgType == websocket.CloseMessage {
			fmt.Printf("conn netid[%d] read CloseMessage, abort read!\n", netid)
			return
		} else if msgType == websocket.BinaryMessage {
			mgr.OnRecvBinaryMessage(netid, buff)
		} else if msgType == websocket.TextMessage {
			mgr.OnRecvTextMessage(netid, buff)
		}
	}
}

func (mgr *ConnMgr) StartTCPRead(netid int, conn *net.Conn) {
	//TODO 用环形buffer
}

func (mgr *ConnMgr) OnRecvTextMessage(netid int, data []byte) {
	fmt.Printf("conn netid[%d] recv text message, data[%s], len[%d], cap[%d]\n", netid, data, len(data), cap(data))

	MsgQueueMgrInstance().OnRecvWSMsg(netid, data)
}

func (mgr *ConnMgr) OnRecvBinaryMessage(netid int, data []byte) {
	fmt.Printf("conn netid[%d] recv binary message, data[%s], len[%d], cap[%d]\n", netid, data, len(data), cap(data))
}