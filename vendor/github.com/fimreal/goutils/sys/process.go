package sys

import (
	"fmt"
	"os"
	"sort"
	"syscall"

	proc "github.com/shirou/gopsutil/process"
)

// Send syscall.SIGHUP to process(pid)
func ProcessReload(pid int) error {
	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}

	return process.Signal(syscall.SIGHUP)
}

func GetPidByName(processName string) ([]int32, error) {
	processes, err := proc.Processes()
	if err != nil {
		return nil, err
	}

	var pids []int32
	for _, process := range processes {
		if name, err := process.Name(); err == nil && name == processName {
			pids = append(pids, process.Pid)
		}
	}

	return pids, nil
}

/*
取到类似 nginx master/worker 工作模式 master 的 pid 。
原理是找到进程 parent pid 值最小的 pid，则认为是 master pid。
*/
func GetMasterPidByName(processName string) (int, error) {
	processes, err := proc.Processes()
	if err != nil {
		return 1, err
	}

	pids := map[int]int{}
	for _, process := range processes {
		name, err := process.Name()
		if err != nil {
			return 1, err
		}

		if name == processName {
			ppid, err := process.Ppid()
			if err != nil {
				return 1, err
			}
			pids[int(ppid)] = int(process.Pid)
		}
	}

	if len(pids) == 0 {
		return 1, fmt.Errorf("not found process named %s ", processName)
	}

	var p []int

	for ppid := range pids {
		p = append(p, ppid)
	}
	sort.Ints(p)
	return pids[p[0]], nil

}
