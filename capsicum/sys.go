package capsicum

import (
	"os"
	"path"
	"runtime"
	"syscall"

	"golang.org/x/sys/unix"
)

const supportsCloseOnExec bool = true // True for any system that supports Capsicum (so far)

// FIXME: unexported function from os
// syscallMode returns the syscall-specific mode bits from Go's portable mode bits.
func syscallMode(i os.FileMode) (o uint32) {
	o |= uint32(i.Perm())
	if i&os.ModeSetuid != 0 {
		o |= syscall.S_ISUID
	}
	if i&os.ModeSetgid != 0 {
		o |= syscall.S_ISGID
	}
	if i&os.ModeSticky != 0 {
		o |= syscall.S_ISVTX
	}
	// No mapping for Go's ModeTemporary (plan9 only).
	return
}

// OpenFileAt is the generalized open call; most users will use Open
// or Create instead. It opens the named file with specified flag
// (O_RDONLY etc.) and perm, (0666 etc.) if applicable. If successful,
// methods on the returned File can be used for I/O.
// If there is an error, it will be of type *PathError.
// FIXME: this is mostly just os.OpenFile - probably that should call this
func OpenFileAt(f *os.File, name string, flag int, perm os.FileMode) (*os.File, error) {
	chmod := false
	if !supportsCreateWithStickyBit && flag&os.O_CREATE != 0 && perm&os.ModeSticky != 0 {
		if _, err := os.Stat(name); os.IsNotExist(err) {
			chmod = true
		}
	}

	var r int
	for {
		var e error
		r, e = unix.Openat(int(f.Fd()), name, flag|unix.O_CLOEXEC, syscallMode(perm))
		if e == nil {
			break
		}

		// On OS X, sigaction(2) doesn't guarantee that SA_RESTART will cause
		// open(2) to be restarted for regular files. This is easy to reproduce on
		// fuse file systems (see http://golang.org/issue/11180).
		if runtime.GOOS == "darwin" && e == syscall.EINTR {
			continue
		}

		return nil, &os.PathError{"open", name, e}
	}

	// open(2) itself won't handle the sticky bit on *BSD and Solaris
	if chmod {
		// FIXME: change to FchmodAt()
		os.Chmod(name, perm)
	}

	// There's a race here with fork/exec, which we are
	// content to live with. See ../syscall/exec_unix.go.
	if !supportsCloseOnExec {
		syscall.CloseOnExec(r)
	}

	return os.NewFile(uintptr(r), path.Join(f.Name(), name)), nil
}
