package scanning

import (
	"context"
	"fmt"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"golang.org/x/sync/syncmap"

	"golang.org/x/sync/semaphore"

	"github.com/ediblesushi/helios/pkg/config"
	"github.com/ediblesushi/helios/pkg/printing"
)

// PortScanner struct
type PortScanner struct {
	ip   string
	lock *semaphore.Weighted
}

// Ulimit gets the limit for concurrent scans
func Ulimit() int64 {
	var rLimit syscall.Rlimit
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		printing.HeliosLog("ERROR", "Error getting rlimit. Using default 1024")
		return 1024
	}

	return int64(rLimit.Cur)
}

// ScanPort will scan each port of the target
func ScanPort(ip string, port int, timeout time.Duration) bool {
	target := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", target, timeout)

	if err != nil {
		if strings.Contains(err.Error(), "too many open files") {
			time.Sleep(timeout)
			ScanPort(ip, port, timeout)
		}
		return false
	}

	conn.Close()
	return true
}

// Start will start the port scan
func (ps *PortScanner) Start(ports []int, timeout time.Duration, optionsBool map[string]bool) {
	openports := syncmap.Map{}
	wg := sync.WaitGroup{}

	for _, port := range ports {
		ps.lock.Acquire(context.TODO(), 1)
		wg.Add(1)
		go func(port int) {
			defer ps.lock.Release(1)
			defer wg.Done()
			openports.Store(port, ScanPort(ps.ip, port, timeout))
		}(port)
	}

	wg.Wait()

	// Sort results
	keys := make([]int, 0)
	openports.Range(func(key, value interface{}) bool {
		k, ok := key.(int)
		if !ok {
			// this will break iteration
			return false
		}

		keys = append(keys, k)

		return true
	})
	sort.Ints(keys)

	for _, k := range keys {
		open, ok := openports.Load(k)
		if !ok {
			printing.HeliosLog("ERROR", "A bad thing happened")
			os.Exit(0)
		}
		if open == true {
			printing.HeliosLog("OPEN", strconv.Itoa(k)+" open")
		} else {
			if optionsBool["verbose"] == true {
				printing.HeliosLog("CLOSED", strconv.Itoa(k)+" closed")
			}
		}
	}
}

// Scan is main scanning function
func Scan(optionsStr map[string]string, optionsBool map[string]bool) {

	printing.HeliosLog("SYSTEM", "Target: "+optionsStr["target"])
	if optionsStr["ports"] == "" {
		optionsStr["ports"] = config.DEFAULTPORTS
	} else {
		printing.HeliosLog("SYSTEM", "Ports: "+optionsStr["ports"])
	}
	s := strings.Split(optionsStr["ports"], ",")
	p := []int{}
	for _, i := range s {
		if strings.Contains(i, "-") {
			portrange := strings.Split(i, "-")
			f, err := strconv.Atoi(portrange[0])
			if err != nil {
				printing.HeliosLog("ERROR", "Error trying to parse the ports")
				os.Exit(0)
			}
			l, err := strconv.Atoi(portrange[1])
			if err != nil {
				printing.HeliosLog("ERROR", "Error trying to parse the ports")
				os.Exit(0)
			}
			a := make([]int, l-f+1)
			for j := range a {
				p = append(p, f+j)
			}
		} else {
			j, err := strconv.Atoi(i)
			if err != nil {
				printing.HeliosLog("ERROR", "Error trying to parse the ports")
				os.Exit(0)
			}
			p = append(p, j)
		}
	}
	ps := &PortScanner{
		ip:   optionsStr["target"],
		lock: semaphore.NewWeighted(Ulimit()),
	}
	ps.Start(p, 10000*time.Millisecond, optionsBool)
}
