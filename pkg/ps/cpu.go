package ps

import "github.com/shirou/gopsutil/v3/cpu"

// Cpu show cpu info
func Cpu() (err error) {
	cpu.Info()
	return
}
