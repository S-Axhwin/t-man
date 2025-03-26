package main

import (
	"fmt"
	"log"
	"sort"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
)

// Function to get CPU/memory usage color gradient
func getGradientColor(usage float64) ui.Color {
	switch {
	case usage < 40:
		return ui.ColorGreen
	case usage < 70:
		return ui.ColorYellow
	default:
		return ui.ColorRed
	}
}

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	// Memory Usage Gauge
	memBar := widgets.NewGauge()
	memBar.Title = "Memory Usage"
	memBar.LabelStyle = ui.NewStyle(ui.ColorWhite)

	// CPU Usage Gauge
	cpuBar := widgets.NewGauge()
	cpuBar.Title = "CPU Usage"
	cpuBar.LabelStyle = ui.NewStyle(ui.ColorWhite)

	// Process Table
	processTable := widgets.NewTable()
	processTable.Title = "Running Processes"
	processTable.RowSeparator = false
	processTable.TextAlignment = ui.AlignLeft
	processTable.Rows = [][]string{{"PID", "Name", "CPU%", "Memory%"}}
	processTable.SetRect(0, 0, 80, 20)

	// Scroll variables
	scrollOffset := 0
	const maxVisibleRows = 15

	// Layout Grid
	grid := ui.NewGrid()
	tWidth, tHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, tWidth, tHeight)
	grid.Set(
		ui.NewRow(0.3, ui.NewCol(0.5, memBar), ui.NewCol(0.5, cpuBar)),
		ui.NewRow(0.7, processTable),
	)

	ui.Render(grid)

	// Variables to store process data
	var procData [][]string

	// Update system usage every second
	go func() {
		for range time.NewTicker(1 * time.Second).C {
			// Fetch CPU Usage
			cpuPercent, _ := cpu.Percent(0, false)
			if len(cpuPercent) > 0 {
				cpuUsage := cpuPercent[0]
				cpuBar.Percent = int(cpuUsage)
				cpuBar.Label = fmt.Sprintf("%.2f%%", cpuUsage)
				cpuBar.BarColor = getGradientColor(cpuUsage)
			}

			// Fetch Memory Usage
			memStats, _ := mem.VirtualMemory()
			memUsage := memStats.UsedPercent
			memBar.Percent = int(memUsage)
			memBar.Label = fmt.Sprintf("%.2f%% (%.2fGB/%.2fGB)", memUsage, float64(memStats.Used)/1e9, float64(memStats.Total)/1e9)
			memBar.BarColor = getGradientColor(memUsage)

			ui.Render(grid)
		}
	}()

	// Update process list every 3 seconds
	go func() {
		for range time.NewTicker(3 * time.Second).C {
			procs, _ := process.Processes()
			newProcData := [][]string{}

			for _, p := range procs {
				pid := p.Pid
				name, _ := p.Name()
				cpuPercent, _ := p.CPUPercent() // Takes time
				memPercent, _ := p.MemoryPercent()

				newProcData = append(newProcData, []string{
					fmt.Sprintf("%d", pid),
					name,
					fmt.Sprintf("%.2f%%", cpuPercent),
					fmt.Sprintf("%.2f%%", memPercent),
				})
			}

			// Sort processes by CPU usage (descending order)
			sort.Slice(newProcData, func(i, j int) bool {
				return newProcData[i][2] > newProcData[j][2]
			})

			// Keep only the top 50 processes
			if len(newProcData) > 50 {
				newProcData = newProcData[:50]
			}

			procData = newProcData // Update global process data
		}
	}()

	// Handle user input (exit on 'q', scroll with up/down)
	for e := range ui.PollEvents() {
		if e.Type == ui.KeyboardEvent {
			switch e.ID {
			case "q":
				return
			case "<Up>":
				if scrollOffset > 0 {
					scrollOffset--
				}
			case "<Down>":
				if scrollOffset+maxVisibleRows < len(procData) {
					scrollOffset++
				}
			}

			// Update table rows with new scroll position
			endIndex := scrollOffset + maxVisibleRows
			if endIndex > len(procData) {
				endIndex = len(procData)
			}
			displayRows := append([][]string{{"PID", "Name", "CPU%", "Memory%"}}, procData[scrollOffset:endIndex]...)

			processTable.Rows = displayRows
			ui.Render(grid)
		}
	}
}
