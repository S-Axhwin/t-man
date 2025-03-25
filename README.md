# Tman - Terminal Process Monitor

Tman is a terminal-based process monitoring tool written in Go. It provides a real-time dashboard displaying CPU and memory usage, along with a list of active processes sorted by CPU usage.

## Features

- Real-time CPU and memory usage monitoring
- Live process tracking sorted by CPU usage
- Scrollable process list
- Keyboard shortcuts for easy navigation

## Installation

Ensure you have Go installed. Clone the repository and build the project:

```sh
git clone https://github.com/your-repo/tman.git
cd tman
go build -o tman
./tman
```

## Dependencies

This project uses the following libraries:

- [gizak/termui](https://github.com/gizak/termui) - For terminal UI components
- [shirou/gopsutil](https://github.com/shirou/gopsutil) - For fetching system metrics

Install dependencies using:

```sh
go get github.com/gizak/termui/v3 github.com/shirou/gopsutil
```

## Code Breakdown

### Importing Required Packages

```go
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
```

- **termui**: Handles the terminal-based UI.
- **gopsutil**: Fetches system resource data like CPU, memory, and process stats.

### UI Elements

#### Memory and CPU Usage Gauges

```go
memBar := widgets.NewGauge()
memBar.Title = "Memory Usage"
memBar.LabelStyle = ui.NewStyle(ui.ColorWhite)

cpuBar := widgets.NewGauge()
cpuBar.Title = "CPU Usage"
cpuBar.LabelStyle = ui.NewStyle(ui.ColorWhite)
```

- These gauges visually represent the memory and CPU usage in real-time.

#### Process Table

```go
processTable := widgets.NewTable()
processTable.Title = "Running Processes"
processTable.Rows = [][]string{{"PID", "Name", "CPU%", "Memory%"}}
```

- Displays process data with columns for Process ID, Name, CPU%, and Memory%.

### Fetching System Metrics

#### CPU and Memory Usage

```go
cpuPercent, _ := cpu.Percent(time.Second, false)
memStats, _ := mem.VirtualMemory()
```

- Retrieves the current CPU and memory usage using `gopsutil`.

#### Fetching Running Processes

```go
procs, _ := process.Processes()
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
```

- Retrieves a list of running processes and extracts CPU and memory usage for each.

### Sorting and Displaying Processes

```go
sort.Slice(procData, func(i, j int) bool {
    return procData[i][2] > procData[j][2] // Sort by CPU usage
})
```

- Ensures the most CPU-intensive processes appear first in the list.

### Keyboard Controls

```go
for e := range ui.PollEvents() {
    switch e.ID {
    case "<C-c>", "q":
        return // Quit on 'q'
    case "<Up>":
        if scrollOffset > 0 { scrollOffset-- }
    case "<Down>":
        if scrollOffset+maxVisibleRows < len(procData) { scrollOffset++ }
    }
    ui.Render(grid)
}
```

- **`q or (crl + c)`**: Quit the program
- **Up/Down Arrows**: Scroll through process list

## Usage

Run the compiled binary:

```sh
./tman
```

Press `q` or `crl + c` to exit and use the arrow keys to navigate the process list.
# t-man
