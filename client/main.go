package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
	"io"
	"log"
	"net/http"
	"time"
)

type DiskUsage struct {
	Device     string
	Mountpoint string
	Total      uint64
	Used       uint64
	Free       uint64
	Percent    float64
}

type Device struct {
	Name      string
	Uptime    uint64
	LocalTime time.Time
	OS        string
	Platform  string
	CPUUsage  float64
	RAM       uint64
	RAMInUse  uint64
	RAMUsage  float64
	Disks     []DiskUsage
	Battery   int
	CPUInfo   []cpu.InfoStat
	NetIO     []net.IOCountersStat
	LoadAvg   *load.AvgStat
	Processes []*process.Process
	// Add more fields as needed
}

func main() {
	// Sample device
	device := Device{}

	// Get memory info
	memoryStat, _ := mem.VirtualMemory()
	device.RAM = memoryStat.Total / 1024 / 1024
	device.RAMInUse = memoryStat.Used / 1024 / 1024
	device.RAMUsage = memoryStat.UsedPercent

	// Get network info
	netStat, _ := net.IOCounters(true)
	device.NetIO = netStat

	// Get disk info
	partitions, _ := disk.Partitions(false)
	for _, partition := range partitions {
		usageStat, _ := disk.Usage(partition.Mountpoint)
		device.Disks = append(device.Disks, DiskUsage{
			Device:     partition.Device,
			Mountpoint: partition.Mountpoint,
			Total:      usageStat.Total / 1024 / 1024,
			Used:       usageStat.Used / 1024 / 1024,
			Free:       usageStat.Free / 1024 / 1024,
			Percent:    usageStat.UsedPercent,
		})
	}

	// Get host info
	hostStat, _ := host.Info()
	device.OS = hostStat.OS
	device.Platform = hostStat.Platform
	device.Name = hostStat.Hostname
	device.Uptime = hostStat.Uptime
	device.LocalTime = time.Now()

	// Get CPU info
	cpuStat, _ := cpu.Info()
	device.CPUInfo = cpuStat
	cpuUsage, _ := cpu.Percent(0, false)
	if len(cpuUsage) > 0 {
		device.CPUUsage = cpuUsage[0]
	}
	if device.OS != "windows" {
		loadAvg, _ := load.Avg()
		device.LoadAvg = loadAvg
	} else {
		device.LoadAvg = nil
	}
	
	jsonDevice, err := json.MarshalIndent(device, "", "  ")
	if err != nil {
		fmt.Println(err)
	}

	resp, err := http.Post("http://localhost:5000/upload_stats", "application/json", bytes.NewBuffer(jsonDevice))
	if err != nil {
		log.Fatalln(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Error closing response body:", err)
		}
	}(resp.Body)

	if resp.StatusCode != 200 {
		log.Fatalln("Error: Status code is not 200")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(body))
}
