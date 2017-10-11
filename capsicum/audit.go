package capsicum

// #include "netinet/tcp.h"
import "C"

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path"
	"strconv"
	"strings"
	"syscall"
)

const (
	fdDir = "/proc/self/fd"
	tcp6  = "/proc/self/net/tcp6"
)

type handler map[string]func(*syscall.Stat_t, string) (FDInfo, error)

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

func null(_ *syscall.Stat_t, _ string) (FDInfo, error) {
	return &NilInfo{}, nil
}

func parseIP6(ips string) (net.IP, uint16, error) {
	if len(ips) != 37 {
		return nil, 0, errorf("Bad IPv6 format '%s'", ips)
	}

	ip := make([]byte, 16)
	for n := 0; n < 16; n++ {
		t, err := strconv.ParseUint(ips[n*2+1:n*2+3], 16, 8)
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
)

type Address struct {
	ip   net.IP
	port uint16
}

type FDSocket struct {
	inode  int
	status SocketStatus
	src    Address
}

func (i FDSocket) String() string {
	if i.status != LISTEN {
		panic("bad status")
	}
	return fmt.Sprintf("LISTEN(%s:%d)", i.src.ip, i.src.port)
}

func listSockInner(f []string, s *FDSocket) error {
	status, err := strconv.ParseInt(f[3], 16, 8)
	if err != nil {
		return err
	}
	if status != C.TCP_LISTEN {
		return errors.New(fmt.Sprintf("Don't know status %d", status))
	}
	ip, port, err := parseIP6(f[1])
	if err != nil {
		return err
	}
	fmt.Printf(" LISTEN(%s:%d)", ip, port)
	s.status = LISTEN
	s.src.ip = ip
	s.src.port = port
	return nil
}

func listSock(_ *syscall.Stat_t, s string) (FDInfo, error) {
	if s[0] != '[' || s[len(s)-1] != ']' {
		return nil, errorf("Can't parse '%s'", s)
	}
	inode, err := strconv.Atoi(s[1 : len(s)-1])
	if err != nil {
		return nil, err
	}
	//fmt.Printf(" inode %d", inode)
	f, err := os.Open(tcp6)
	defer f.Close()
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
			panic(err)
		}
		if i != inode {
			continue
		}
		s := FDSocket{inode: inode}
		//fmt.Printf(" %#v", f)
		listSockInner(f, &s)
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
	files, err := ioutil.ReadDir(fdDir)
	if err != nil {
		return nil, err
	}
	fds := make([]*FD, 0)
	for _, file := range files {
		var fd FD
		fds = append(fds, &fd)

		i, err := strconv.Atoi(file.Name())
		if err != nil {
			return nil, err
		}
		if i < 0 {
			return nil, errorf("FD %d is negative", i)
		}
		fd.fd = uintptr(i)

		bname := path.Join(fdDir, file.Name())
		name, err := os.Readlink(bname)
		if err != nil {
			if err.(*os.PathError).Err == syscall.ENOENT {
				fmt.Printf("%s disappeared\n", file.Name())
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
				i, err := os.Lstat(bname)
				if err != nil {
					return nil, err
				} else {
					stat := i.Sys().(*syscall.Stat_t)
					info, err := f(stat, scheme[1])
					if err != nil {
						return nil, err
					}
					fd.info = info
				}
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
