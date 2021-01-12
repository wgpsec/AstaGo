package main

import (
	"AstaGo/Tools/PortScan"
	"flag"
	"fmt"
	"strings"
	"sync"
	"time"
)



// ----------------------------------------------------------------------------
//端口扫描测试
func main() {
	var scanports []string
	defaultports := [...]string{"21", "22", "23", "25", "80", "443", "8080",
		"110", "135", "139", "445", "389", "489", "587", "1433", "1434",
		"1521", "1522", "1723", "2121", "3306", "3389", "4899", "5631",
		"5632", "5800", "5900", "7071", "43958", "65500", "4444", "8888",
		"6789", "4848", "5985", "5986", "8081", "8089", "8443", "10000",
		"6379", "7001", "7002"}
	var ip string
	var ports string
	var gonum int
	flag.StringVar(&ip, "i", "127.0.0.1", "扫描的IP地址，默认本机地址")
	flag.StringVar(&ports, "p", "", "扫描的端口地址，默认使用默认端口")
	flag.IntVar(&gonum, "g", 4, "开启的goroutine的数量，默认为4")
	flag.StringVar(&taskfile, "r", "", "待扫描的的IP和对应的端口号所在的文件")
	flag.Parse()


	if len(ports) != 0 {
		scanports = strings.Split(ports, ",")
	} else {
		scanports = defaultports[:]
	}
	taskschan := make(chan string, len(scanports))
	reschan := make(chan string, len(scanports))
	exitchan := make(chan bool, 4)
	var wgp sync.WaitGroup
	for _, value := range scanports {
		taskschan <- value
	}
	close(taskschan)
	start := time.Now()
	for i := 0; i < gonum; i++ {
		wgp.Add(1)
		go PortScan.Scan("127.0.0.1", taskschan, reschan, exitchan, &wgp)
	}
	wgp.Wait()
	for i := 0; i < gonum; i++ {
		<-exitchan
	}
	end := time.Since(start)
	close(exitchan)
	close(reschan)
	for {
		openport, ok := <-reschan
		if !ok {
			break
		}
		fmt.Println("开放的端口：", openport)
	}
	fmt.Println("花费的时间：", end)
}

