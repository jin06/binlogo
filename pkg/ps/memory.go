package ps

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/sirupsen/logrus"
)

func Memory() (err error){
	v, err := mem.VirtualMemory()
	if err != nil {
		return
	}
	logrus.Debug(v.Total, v.Free, v.UsedPercent)
	fmt.Println("total", v.Total)
	fmt.Println("free", v.Free)
	fmt.Println("usedPercent", v.UsedPercent)
	return
}
