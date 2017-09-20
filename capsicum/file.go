package capsicum

import "os"

type hasfd interface {
	Fd() uintptr
}

func CapRightsLimit(f *os.File, r *CapRights) error {
	return CapRightsLimitFd(f.Fd(), r)
}

func CapRightsGet(f *os.File) (*CapRights, error) {
	return CapRightsGetFd(f.Fd())
}
