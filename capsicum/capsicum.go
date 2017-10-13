package capsicum

import (
	"os"
	"sort"
	"syscall"
)

// like ioutil.ReadDir, but from an already open file
func ReaddirnamesAt(f *os.File) ([]string, error) {
	list, err := f.Readdirnames(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	sort.Slice(list, func(i, j int) bool { return list[i] < list[j] })
	return list, nil
}

func Dup(f *os.File) (*os.File, error) {
	fd, err := syscall.Dup(int(f.Fd()))
	if err != nil {
		return nil, err
	}
	return os.NewFile(uintptr(fd), f.Name()), nil
}
