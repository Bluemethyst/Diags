package main

import (
	"encoding/json"
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
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
	OS        string
	Platform  string
	RAM       int
	RamInUse  int
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
	device.RAM = int(memoryStat.Total / 1024 / 1024)
	device.RamInUse = int(memoryStat.Used / 1024 / 1024)

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

	// Get CPU info
	cpuStat, _ := cpu.Info()
	device.CPUInfo = cpuStat

	// Get load average
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

	fmt.Println(string(jsonDevice))
}
