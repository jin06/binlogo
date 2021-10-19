package ps

import "github.com/shirou/gopsutil/v3/cpu"

func Cpu() (err error){
	cpu.Info()
	return
}
