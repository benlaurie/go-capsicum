package capsicum

// #include "netinet/tcp.h"
import "C"

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"syscall"
)

const (
	fdDir = "/proc/self/fd"
	tcp6  = "/proc/self/net/tcp6"
)

var fdDirFile, tcp6File *os.File

func openOrDie(path string) *os.File {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	return f
}

func init() {
	fdDirFile = openOrDie(fdDir)
	tcp6File = openOrDie(tcp6)
}

type handler map[string]func(string) (FDInfo, error)

var handlers = handler{
	"socket":     listSock,
	"anon_inode": null,
}

func errorf(f string, a ...interface{}) error {
	return errors.New(fmt.Sprintf(f, a))
}

type NilInfo struct{}

func (NilInfo) String() string {
	return "[no info]"
}

func null(_ string) (FDInfo, error) {
	return &NilInfo{}, nil
}

func parseIP6(ips string) (net.IP, uint16, error) {
	if len(ips) != 37 {
		return nil, 0, errorf("Bad IPv6 format '%s'", ips)
	}

	ip := make([]byte, 16)
	for n := 0; n < 16; n++ {
		t, err := strconv.ParseUint(ips[n*2:n*2+2], 16, 8)
		if err != nil {
			return nil, 0, err
		}
		ip[n] = byte(t)
	}
	p, err := strconv.ParseUint(ips[33:37], 16, 16)
	if err != nil {
		return nil, 0, err
	}
	return ip, uint16(p), nil
}

type SocketStatus int

const (
	LISTEN SocketStatus = iota
	CLIENT
	UNKNOWN
)

type Address struct {
	ip   net.IP
	port uint16
}

type FDSocket struct {
	inode         int
	status        SocketStatus
	unknownStatus uint64 // The raw status from the socket - only set if |status| is UNKNOWN
	src           Address
	dest          Address // Not set if |status| is LISTEN
}

func (i FDSocket) String() string {
	switch i.status {
	case LISTEN:
		return fmt.Sprintf("LISTEN(%s:%d)", i.src.ip, i.src.port)
	case CLIENT:
		return fmt.Sprintf("CLIENT(%s:%d -> %s:%d)", i.src.ip, i.src.port, i.dest.ip, i.dest.port)
	case UNKNOWN:
		return fmt.Sprintf("SOCKET[%d](%s:%d -> %s:%d)", i.unknownStatus, i.src.ip, i.src.port, i.dest.ip, i.dest.port)
	default:
		panic("Unknwon status")
	}
}

func listSockInner(f []string, s *FDSocket) error {
	status, err := strconv.ParseUint(f[3], 16, 8)
	if err != nil {
		return err
	}
	ip, port, err := parseIP6(f[1])
	if err != nil {
		return err
	}
	s.src.ip = ip
	s.src.port = port

	if status == C.TCP_LISTEN {
		s.status = LISTEN
		return nil
	}

	ip, port, err = parseIP6(f[2])
	if err != nil {
		return err
	}
	s.dest.ip = ip
	s.dest.port = port

	if status == C.TCP_ESTABLISHED {
		s.status = CLIENT
		return nil
	}

	s.status = UNKNOWN
	s.unknownStatus = status

	return nil
}

func listSock(s string) (FDInfo, error) {
	if s[0] != '[' || s[len(s)-1] != ']' {
		return nil, errorf("Can't parse '%s'", s)
	}
	inode, err := strconv.Atoi(s[1 : len(s)-1])
	if err != nil {
		return nil, err
	}
	//fmt.Printf(" inode %d", inode)
	f := tcp6File
	_, err = f.Seek(0, os.SEEK_SET)
	if err != nil {
		return nil, err
	}
	r := bufio.NewReader(f)
	l, p, err := r.ReadLine()
	if err != nil {
		return nil, err
	}
	if p {
		return nil, errors.New("line too long")
	}
	if string(l) != "  sl  local_address                         remote_address                        st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode" {
		return nil, errorf("unknown format: %s", l)
	}
	for l, p, err = r.ReadLine(); ; l, p, err = r.ReadLine() {
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}
		f := strings.Fields(string(l))
		if len(f) < 11 {
			errorf("Don't understand: %s", l)
		}
		i, err := strconv.Atoi(f[9])
		if err != nil {
			return nil, err
		}
		if i != inode {
			continue
		}
		s := FDSocket{inode: inode}
		//fmt.Printf(" %#v", f)
		err = listSockInner(f, &s)
		if err != nil {
			return nil, err
		}
		return &s, nil
	}
	return nil, errorf("socket %d not found", inode)
}

type FDType int

const (
	FILE FDType = iota
	SOCKET
)

type FDInfo interface {
	String() string
}

type FD struct {
	fd     uintptr
	name   string
	rights *CapRights
	info   FDInfo
}

func (fd FD) String() string {
	return fmt.Sprintf("%d -> %s %s %s", fd.fd, fd.name, fd.info, fd.rights)
}

type FDDisappeared struct{}

func (FDDisappeared) String() string {
	return "disappeared"
}

type FDFile struct{}

func (FDFile) String() string {
	return "FILE"
}

// FIXME: an evil program could mess with this by dup()ing and close()ing a lot...
func GetAllFDInfo() ([]*FD, error) {
	// ReaddirnamesAt consumes its File
	dir, err := Dup(fdDirFile)
	if err != nil {
		return nil, err
	}

	files, err := ReaddirnamesAt(dir)
	if err != nil {
		return nil, err
	}

	fds := make([]*FD, 0)
	for _, file := range files {
		var fd FD
		fds = append(fds, &fd)

		i, err := strconv.Atoi(file)
		if err != nil {
			return nil, err
		}
		if i < 0 {
			return nil, errorf("FD %d is negative", i)
		}
		fd.fd = uintptr(i)

		//bname := path.Join(fdDir, file.Name())
		//name, err := os.Readlink(bname)
		name, err := ReadlinkAt(fdDirFile, file)
		if err != nil {
			if err == syscall.ENOENT {
				fd.info = FDDisappeared{}
			} else {
				return nil, err
			}
			continue
		}
		fd.name = name

		// FIXME: if fd has disappeared, this would panic, hence run after disappearance check, but unless we lock all other threads, there's a race...
		r, err := CapRightsGetFd(fd.fd)
		if err != nil {
			return nil, err
		}
		fd.rights = r

		// FIXME: filename could have this format, but obvs would then not belong to the scheme...
		scheme := strings.SplitN(name, ":", 2)
		if len(scheme) > 1 {
			f := handlers[scheme[0]]
			if f != nil {
				info, err := f(scheme[1])
				if err != nil {
					return nil, err
				}
				fd.info = info
			} else {
				return nil, errorf(" no handler for '%s'", scheme[0])
			}
		} else {
			fd.info = &FDFile{}
		}
	}
	return fds, nil
}

func ListAllFDs() error {
	fds, err := GetAllFDInfo()
	if err != nil {
		return err
	}
	fmt.Println("---")
	defer fmt.Println("---")
	for _, fd := range fds {
		fmt.Println(fd)
	}
	return nil
}
