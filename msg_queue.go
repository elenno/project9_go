package main

import (
	"container/list"
	"sync"
	"time"
	"fmt"
)

type Msg struct {
	Netid int
	Data []byte
}

type MsgQueue struct {
	cur_size int
	max_size int
	data_holder *list.List
	mutex sync.Mutex
}

func NewMsgQueue() *MsgQueue {
    return &MsgQueue {
		cur_size: 0,
		max_size: 100000,
		data_holder: list.New(),
	}
}

//加锁在外部加
func (q *MsgQueue) isEmpty() bool {
	return q.cur_size <= 0 
}

func (q *MsgQueue) isFull() bool {
	return q.cur_size >= q.max_size
}

func (q *MsgQueue) IsFull() bool {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return q.isFull()
}

func (q *MsgQueue) IsEmpty() bool {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return q.isEmpty()
}

func (q *MsgQueue) Push(value interface{}) bool {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if q.isFull() {
		return false
	}

	q.data_holder.PushBack(value)
	q.cur_size += 1

	return true
}

func (q *MsgQueue) Pop() interface{} {
	q.mutex.Lock()	
	defer q.mutex.Unlock()

	if q.isEmpty() {
		return nil
	}
	
	e := q.data_holder.Front()
	q.data_holder.Remove(e)
	q.cur_size -= 1
	
	return e.Value  
}

func (q *MsgQueue) GetCurSize() int {
	return q.cur_size
}

type MsgQueueMgr struct {
	ws_recv_queue *MsgQueue
	tcp_recv_queue *MsgQueue
}

var msgQueueMgrInstance *MsgQueueMgr
var msgQueueInstMutex sync.Mutex
func MsgQueueMgrInstance() *MsgQueueMgr {
	if nil == msgQueueMgrInstance {
		msgQueueInstMutex.Lock()
		if nil == msgQueueMgrInstance {
			msgQueueMgrInstance = &MsgQueueMgr{
				ws_recv_queue: NewMsgQueue(),
				tcp_recv_queue: NewMsgQueue(),
			}
		}
		msgQueueInstMutex.Unlock()   
	}
    return msgQueueMgrInstance
}

func (mgr *MsgQueueMgr) StartMsgQueue() {
	fmt.Printf("StartMsgQueue\n")

	go mgr.RunWSRecvQueue()
	go mgr.RunTCPRecvQueue()
}

func (mgr *MsgQueueMgr) OnRecvWSMsg(netid int, data []byte) {
	var m Msg
	m.Netid = netid
	m.Data = data
	mgr.ws_recv_queue.Push(m)
}

func (mgr *MsgQueueMgr) OnRecvTCPMsg(netid int, data []byte) {
	var m Msg
	m.Netid = netid
	m.Data = data
	mgr.tcp_recv_queue.Push(m)
}

func (mgr *MsgQueueMgr) RunWSRecvQueue() {
	for {
		frontVal := mgr.ws_recv_queue.Pop()
		if nil == frontVal {
			time.Sleep(1 * time.Microsecond)
			continue
		}

		msg := frontVal.(Msg)
		fmt.Printf("RunWSRecvQueue from netid[%d] data[%s]\n", msg.Netid, msg.Data)

		// TODO 收到信息，发到gameworld
	}
}

func (mgr *MsgQueueMgr) RunTCPRecvQueue() {
	for {
		frontVal := mgr.tcp_recv_queue.Pop()
		if nil == frontVal {
			time.Sleep(1 * time.Microsecond)
			continue
		}

		msg := frontVal.(Msg)
		fmt.Printf("RunTCPRecvQueue from netid[%d] data[%s]\n", msg.Netid, msg.Data)
	}
}