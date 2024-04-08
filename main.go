package main

import (
	"fmt"
	"log"
	"time"

	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/process"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"

	"github.com/fatih/structs"
)

type ProcessAdapter struct {
	Pid        int32   `json:"pid"`
	Name       string  `json:"name"`
	NumThreads int32   `json:"num_threads"`
	RSSMem     uint64  `json:"rss_mem"`
	CPUPercent float64 `json:"cpu_percent"`
}

func convertProcess(inputProcess *process.Process, processAdapter *ProcessAdapter) {
	processAdapter.Pid = inputProcess.Pid
	processAdapter.Name, _ = inputProcess.Name()
	processAdapter.NumThreads, _ = inputProcess.NumThreads()
	meminfo, _ := inputProcess.MemoryInfo()
	processAdapter.RSSMem = meminfo.RSS

	start := time.Now()
	processAdapter.CPUPercent, _ = inputProcess.CPUPercent()
	elapsed := time.Since(start)
	log.Printf("Get cpu took %s", elapsed)
}

var data = [][]string{[]string{"top left", "top right"},
	[]string{"bottom left", "bottom right"}}

func main() {
	infoStat, _ := host.Info()
	fmt.Printf("Total processes: %d\n", infoStat.Procs)

	miscStat, _ := load.Misc()
	fmt.Printf("Running processes: %d\n", miscStat.ProcsRunning)

	start := time.Now()
	processes, _ := process.Processes()
	elapsed := time.Since(start)
	log.Printf("Get processes took %s", elapsed)

	start = time.Now()
	processAdapters := make([]*ProcessAdapter, len(processes))
	for i := range len(processes) {
		processAdapters[i] = new(ProcessAdapter)
	}
	for i, p := range processes {
		start1 := time.Now()
		convertProcess(p, processAdapters[i])
		elapsed1 := time.Since(start1)
		log.Printf("Convert 1 process took %s", elapsed1)
	}
	elapsed = time.Since(start)
	log.Printf("Convert processes took %s", elapsed)

	myApp := app.New()
	myWindow := myApp.NewWindow("Table Widget")

	start = time.Now()
	list := widget.NewTable(
		func() (int, int) {
			return len(processAdapters), 5
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("wide content")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(fmt.Sprint(structs.Values(processAdapters[i.Row])[i.Col]))
		})
	elapsed = time.Since(start)
	log.Printf("Create table took %s", elapsed)

	myWindow.SetContent(list)
	myWindow.ShowAndRun()
}
