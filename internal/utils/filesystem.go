package utils

import (
	"fmt"
	"log"
	"strconv"
	"syscall"
)

const (
	AFS          = 26
	EXT4         = 61267
	TMPFS        = 16914836
	XFS          = 1481003842
	NTFS_FUSEBLK = 1702057286
	BTRFS        = 2435016766
)

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

type DiskStatus struct {
	All   uint64 `json:"all"`
	Used  uint64 `json:"used"`
	Free  uint64 `json:"free"`
	Avail uint64 `json:"avail"`
}

func CheckFS(path string) string {

	var fs syscall.Statfs_t

	err := syscall.Statfs(path, &fs)
	if err != nil {
		log.Println(err)
	}
	getDiskSpace(fs, "GB")

	return getFsTypeName(fs.Type)
}

func getFsTypeName(fsType int64) string {
	switch fsType {
	case EXT4:
		return "ext4"
	case TMPFS:
		return "tmpfs"
	case XFS:
		return "xfs"
	case BTRFS:
		return "btrfs"
	case NTFS_FUSEBLK:
		return "ntfs"
	case AFS:
		return "afs"
	default:
		return "can`t find fsname " + strconv.FormatUint(uint64(fsType), 10)
	}
}

func getDiskSpace(fs syscall.Statfs_t, sizetype string) (disk DiskStatus) {

	disk.All = fs.Blocks * uint64(fs.Bsize)
	disk.Avail = fs.Bavail * uint64(fs.Bsize)
	disk.Free = fs.Bfree * uint64(fs.Bsize)
	disk.Used = disk.All - disk.Free

	switch sizetype {
	case "GB":
		// fmt.Printf("Avail: %.2f GB\n", float64(disk.Avail)/float64(GB))
	case "MB":
		// fmt.Printf("Avail: %.2f MB\n", float64(disk.Avail)/float64(MB))
	default:
		fmt.Printf("All: %.2f GB\n", float64(disk.All)/float64(GB))
		fmt.Printf("Used: %.2f GB\n", float64(disk.Used)/float64(GB))
		fmt.Println("")
	}
	return
}
