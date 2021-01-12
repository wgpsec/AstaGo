package PortScan

import (
	"net"
	"sync"
	"time"
)

func Scan(ip string, taskschan chan string, reschan chan string, exitchan chan bool, wgscan *sync.WaitGroup) {
	defer func() {
		exitchan <- true
		wgscan.Done()
	}()
	for {
		port, ok := <-taskschan
		if !ok {
			break
		}
		_, err := net.DialTimeout("tcp", ip+":"+port, time.Second)
		if err == nil {
			reschan <- port
		}
	}
}
