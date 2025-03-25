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

// Function to get row background color (alternating gradient effect)
func getRowBgColor(index int) ui.Color {
	colors := []ui.Color{ui.ColorBlue, ui.ColorCyan} // Alternate between two colors
	return colors[index%len(colors)]
}

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
	processTable.SetRect(0, 0, 80, 20) // Set dimensions

	// Scroll variables
	scrollOffset := 0
	const maxVisibleRows = 10 // Number of rows visible at a time

	// Layout Grid
	grid := ui.NewGrid()
	tWidth, tHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, tWidth, tHeight)
	grid.Set(
		ui.NewRow(0.3, ui.NewCol(0.5, memBar), ui.NewCol(0.5, cpuBar)),
		ui.NewRow(0.7, processTable),
	)

	ui.Render(grid)

	ticker := time.NewTicker(time.Second).C
	var procData [][]string

	go func() {
		for range ticker {
			// Fetch CPU Usage
			cpuPercent, _ := cpu.Percent(time.Second, false)
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

			// Fetch Running Processes
			procs, _ := process.Processes()
			procData = [][]string{}

			for _, p := range procs {
				pid := p.Pid
				name, _ := p.Name()
				cpuPercent, _ := p.CPUPercent()
				memPercent, _ := p.MemoryPercent()
				procData = append(procData, []string{
					fmt.Sprintf("%d", pid),
					name,
					fmt.Sprintf("%.2f%%", cpuPercent),
					fmt.Sprintf("%.2f%%", memPercent),
				})
			}

			// Sort processes by CPU usage (descending order)
			sort.Slice(procData, func(i, j int) bool {
				return procData[i][2] > procData[j][2]
			})

			// Keep only the top 30 processes
			if len(procData) > 30 {
				procData = procData[:30]
			}

			// Update table rows with scroll offset
			endIndex := scrollOffset + maxVisibleRows
			if endIndex > len(procData) {
				endIndex = len(procData)
			}
			displayRows := append([][]string{{"PID", "Name", "CPU%", "Memory%"}}, procData[scrollOffset:endIndex]...)

			processTable.Rows = displayRows
			for i := 1; i < len(displayRows); i++ {
				processTable.RowStyles[i] = ui.NewStyle(ui.ColorWhite, getRowBgColor(i)) // Full row color
			}

			ui.Render(grid)
		}
	}()

	// Handle user input (exit on 'q', scroll with up/down)
	for e := range ui.PollEvents() {
		if e.Type == ui.KeyboardEvent {
			switch e.ID {
			case "<C-c>", "q":
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
			ui.Render(grid)
		}
	}
}
