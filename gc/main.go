package main

import (
	"cowmap"
	"fmt"
	"runtime"
	"time"
)

const ItemCount = 20000

func main() {
	ReportMemStats()
	cowmap := cowmap.New()

	fmt.Printf("ItemCount = %d\n", ItemCount)
	for i := 0; i < ItemCount; i++ {
		cowmap.Set(i, i)
	}

	ReportMemStats()

	fmt.Println("Run GC")
	runtime.GC()

	ReportMemStats()

	time.Sleep(time.Hour)
}

func ReportMemStats() {
	memStats := runtime.MemStats{}
	runtime.ReadMemStats(&memStats)

	fmt.Printf("heap alloc: %d bytes\n", memStats.HeapAlloc)

	unit := "KB"
	allocatedMem := memStats.HeapAlloc / 1024

	if allocatedMem >= 1024 {
		allocatedMem /= 1024
		unit = "MB"
	}

	fmt.Printf("heap alloc: %d %s\n", allocatedMem, unit)
}
