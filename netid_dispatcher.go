package main

import (
	"sync"
)

var mutex sync.Mutex
var gNetIDGenerator = 0
const MaxNetID int = 10000000 

func GetNewNetID() int {
	mutex.Lock()
	defer mutex.Unlock()

	gNetIDGenerator = gNetIDGenerator % MaxNetID + 1
	return gNetIDGenerator
}