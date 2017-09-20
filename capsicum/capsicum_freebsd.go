// +build freebsd

// Note that we could just use the capsicum_linux.go module for FreeBSD, too, but instead use the built-in support.

package capsicum

import "golang.org/x/sys/unix"

const (
	ECAPMODE    = unix.ECAPMODE
	ENOTCAPABLE = unix.ENOTCAPABLE
)

const (
	CAP_EVENT  = unix.CAP_EVENT
	CAP_LISTEN = unix.CAP_LISTEN
	CAP_LOOKUP = unix.CAP_LOOKUP
	CAP_PDWAIT = unix.CAP_PDWAIT
	CAP_READ   = unix.CAP_READ
	CAP_WRITE  = unix.CAP_WRITE
)

type CapRights unix.CapRights

func CapEnter() error {
	return unix.CapEnter()
}

func CapRightsInit(rights ...uint64) (*CapRights, error) {
	r, err := unix.CapRightsInit(rights)
	return (*CapRights)(r), err
}

// FIXME: should take a File, not an fd?
func CapRightsLimitFd(fd uintptr, r *CapRights) error {
	return unix.CapRightsLimit(fd, (*unix.CapRights)(r))
}

func CapRightsGetFd(fd uintptr) (*CapRights, error) {
	r, err := unix.CapRightsGet(fd)
	return (*CapRights)(r), err
}

func CapRightsSet(r *CapRights, rights ...uint64) error {
	return unix.CapRightsSet((*unix.CapRights)(r), rights)
}

func CapRightsClear(r *CapRights, rights ...uint64) error {
	return unix.CapRightsClear((*unix.CapRights)(r), rights)
}

func CapRightsIsSet(r *CapRights, rights ...uint64) (bool, error) {
	return unix.CapRightsIsSet((*unix.CapRights)(r), rights)
}
