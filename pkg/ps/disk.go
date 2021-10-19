package ps

 import "github.com/shirou/gopsutil/v3/disk"


func Disk() (err error) {
 disk.Partitions()
}
