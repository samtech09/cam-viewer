package main

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	// "github.com/mackerelio/go-osstat/cpu"
	// "github.com/mackerelio/go-osstat/memory"
)

type Stats struct {
	Cpu        string
	Mem        memory
	Swap       memory
	Disk       memory
	DiskPath   string
	RefreshInt int
}

type memory struct {
	Total       string
	Free        string
	Used        string
	UsedPercent float64
}

func GetStats() Stats {
	s := Stats{}
	s.Cpu = getCPU()
	s.Disk = getDisk(appConfig.StatsDiskPath)
	s.DiskPath = appConfig.StatsDiskPath
	s.Mem = getMEM()
	s.Swap = getSWAP()
	s.RefreshInt = appConfig.StatsRefreshInterval
	return s
}

func getCPU() string {
	percent, _ := cpu.Percent(time.Second, false)
	//fmt.Printf("  User: %.2f\n", percent[0])
	return fmt.Sprintf("%.2f", percent[0])
	// fmt.Printf("  User: %.2f\n", percent[cpu.CPUser])
	// fmt.Printf("  Nice: %.2f\n", percent[cpu.CPNice])
	// fmt.Printf("   Sys: %.2f\n", percent[cpu.CPSys])
	// fmt.Printf("  Intr: %.2f\n", percent[cpu.CPIntr])
	// fmt.Printf("  Idle: %.2f\n", percent[cpu.CPIdle])
	// fmt.Printf("States: %.2f\n", percent[cpu.CPUStates])
}

func getMEM() memory {
	v, _ := mem.VirtualMemory()
	m := memory{}
	m.Total = BytesToHuman(v.Total)
	m.Free = BytesToHuman(v.Free)
	m.Used = BytesToHuman(v.Used)
	m.UsedPercent = v.UsedPercent
	// // almost every return value is a struct
	// fmt.Printf("MEM Total: %s, Free:%s, UsedPercent:%f%%\n", m.Total, m.Free, m.UsedPercent)
	return m
}

func getSWAP() memory {
	v, _ := mem.SwapMemory()
	m := memory{}
	m.Total = BytesToHuman(v.Total)
	m.Free = BytesToHuman(v.Free)
	m.Used = BytesToHuman(v.Used)
	m.UsedPercent = v.UsedPercent
	// // almost every return value is a struct
	// fmt.Printf("SWAP Total: %s, Free:%s, UsedPercent:%f%%\n", m.Total, m.Free, m.UsedPercent)
	return m
}

func getDisk(path string) memory {
	v, _ := disk.Usage(path)
	m := memory{}
	m.Total = BytesToHuman(v.Total)
	m.Free = BytesToHuman(v.Free)
	m.Used = BytesToHuman(v.Used)
	m.UsedPercent = v.UsedPercent
	// // almost every return value is a struct
	// fmt.Printf("DISK Total: %s, Free:%s, UsedPercent:%f%%\n", m.Total, m.Free, m.UsedPercent)
	return m
}

// func getCPU() float64 {
// 	before, err := cpu.Get()
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "%s\n", err)
// 		return 0
// 	}
// 	time.Sleep(time.Duration(1) * time.Second)
// 	after, err := cpu.Get()
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "%s\n", err)
// 		return 0
// 	}
// 	total := float64(after.Total - before.Total)
// 	fmt.Printf("cpu user: %f %%\n", float64(after.User-before.User)/total*100)
// 	fmt.Printf("cpu system: %f %%\n", float64(after.System-before.System)/total*100)
// 	fmt.Printf("cpu idle: %f %%\n", float64(after.Idle-before.Idle)/total*100)

// 	return total * 100
// }

// func getMEM() uint64 {
// 	memory, err := memory.Get()
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "%s\n", err)
// 		return 0
// 	}
// 	fmt.Printf("memory total: %d bytes\n", memory.Total)
// 	fmt.Printf("memory used: %d bytes\n", memory.Used)
// 	fmt.Printf("memory cached: %d bytes\n", memory.Cached)
// 	fmt.Printf("memory free: %d bytes\n", memory.Free)

// 	return memory.Total
// }
