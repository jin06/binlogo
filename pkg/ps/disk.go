package ps

import "github.com/shirou/gopsutil/v3/disk"

// Disk show disk info
func Disk() (err error) {
	disk.Partitions(true)
	return
}
