package main

import (
	"encoding/json"
	"fmt"

	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/process"
)

type ProcessAdapter struct {
	Pid        int32   `json:"pid"`
	Name       string  `json:"name"`
	NumThreads int32   `json:"num_threads"`
	RSSMem     uint64  `json:"rss_mem"`
	CPUPercent float64 `json:"cpu_percent"`
}

func convertProcess(inputProcess *process.Process) ProcessAdapter {
	var processAdapter ProcessAdapter

	processAdapter.Pid = inputProcess.Pid
	processAdapter.Name, _ = inputProcess.Name()
	processAdapter.NumThreads, _ = inputProcess.NumThreads()
	meminfo, _ := inputProcess.MemoryInfo()
	processAdapter.RSSMem = meminfo.RSS
	processAdapter.CPUPercent, _ = inputProcess.CPUPercent()

	return processAdapter
}

func main() {
	infoStat, _ := host.Info()
	fmt.Printf("Total processes: %d\n", infoStat.Procs)

	miscStat, _ := load.Misc()
	fmt.Printf("Running processes: %d\n", miscStat.ProcsRunning)

	processes, _ := process.Processes()

	var processAdapter ProcessAdapter
	for _, p := range processes {
		processAdapter = convertProcess(p)
		s, _ := json.Marshal(processAdapter)
		fmt.Printf("%s\n", s)
	}
}
