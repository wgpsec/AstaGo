package IcmpScan

import (
	"io/ioutil"
	"log"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

func ScanHost(ip string, count int) bool {
	var cmd = &exec.Cmd{}
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("ping", "-n", strconv.Itoa(count), ip)
	default:
		cmd = exec.Command("ping", "-c", strconv.Itoa(count), ip)
	}
	output, err := cmd.StdoutPipe()
	defer output.Close()
	cmd.Start()
	if err != nil {
		log.Fatal(err)
	} else {
		result, err := ioutil.ReadAll(output)
		if err != nil {
			log.Fatal(err)
		}
		if strings.Contains(string(result), "TTL") || strings.Contains(string(result), "ttl") {
			return true
		} else {
			return false
		}
	}
	return false
}

func ScanHostTasks(hostsChan chan string, resChan chan string, exitChan chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		host, ok := <-hostsChan
		if !ok {
			break
		} else {
			res := ScanHost(host, 4)
			if res {
				resChan <- host
			}
		}
	}
	exitChan <- true
}
