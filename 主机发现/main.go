package main

import (
	"AstaGo/Tools/IcmpScan"
	"flag"
	"fmt"
	"strings"
	"sync"
	"time"
)



// ----------------------------------------------------------------
//主机发现测试
func main() {
	//hosts := [...]string{
	//	"127.0.0.1", "110.242.68.4", "192.168.10.3"}

	var hostslist string
	var gonum int
	flag.StringVar(&hostslist, "i", "", "输入要扫描的地址，默认为空")
	flag.IntVar(&gonum, "g", 1, "需要开启的goroutine的数量")
	flag.Parse()
	if len(hostslist) == 0 {
		fmt.Println("请输入要扫描的主机")
	} else {
		hosts := strings.Split(hostslist, ",")
		tasksChan := make(chan string, len(hosts))
		resChan := make(chan string, len(hosts))
		exitChan := make(chan bool, 4)
		var wg sync.WaitGroup
		for _, host := range hosts {
			tasksChan <- host
		}
		close(tasksChan)
		start := time.Now()
		for i := 0; i < gonum; i++ {
			wg.Add(1)
			go IcmpScan.ScanHostTasks(tasksChan, resChan, exitChan, &wg)
		}
		wg.Wait()
		for i := 0; i < gonum; i++ {
			<-exitChan
		}
		close(resChan)
		end := time.Since(start)
		for {
			openhost, ok := <-resChan
			if !ok {
				break
			}
			fmt.Println("开放的主机", openhost)
		}
		fmt.Println("花费的时间", end)
	}
}
