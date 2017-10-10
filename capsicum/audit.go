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

type handler map[string]func(*syscall.Stat_t, string) error

var handlers = handler{
	"socket":     listSock,
	"anon_inode": null,
}

func errorf(f string, a ...interface{}) error {
	return errors.New(fmt.Sprintf(f, a))
}

func null(_ *syscall.Stat_t, _ string) error {
	return nil
}

func parseIP6(ips string) (net.IP, uint64, error) {
	if len(ips) != 37 {
		return nil, 0, errorf("Bad IPv6 format '%s'", ips)
	}

	ip := make([]byte, 16)
	for n := 0; n < 16; n++ {
		t, err := strconv.ParseUint(ips[n*2+1:n*2+2], 16, 8)
		if err != nil {
			return nil, 0, err
		}
		ip[n] = byte(t)
	}
	p, err := strconv.ParseUint(ips[33:37], 16, 16)
	if err != nil {
		return nil, 0, err
	}
	return ip, p, nil
}

func listSockInner(f []string) error {
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
	fmt.Printf(" LISTEN(%s:%d)", ip.String(), port)
	return nil
}

func listSock(_ *syscall.Stat_t, s string) error {
	if s[0] != '[' || s[len(s)-1] != ']' {
		return errorf("Can't parse '%s'", s)
	}
	inode, err := strconv.Atoi(s[1 : len(s)-1])
	if err != nil {
		return err
	}
	//fmt.Printf(" inode %d", inode)
	f, err := os.Open(tcp6)
	if err != nil {
		return err
	}
	r := bufio.NewReader(f)
	l, p, err := r.ReadLine()
	if err != nil {
		return err
	}
	if p {
		return errors.New("line too long")
	}
	if string(l) != "  sl  local_address                         remote_address                        st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode" {
		return errorf("unknown format: %s", l)
	}
	for l, p, err = r.ReadLine(); ; l, p, err = r.ReadLine() {
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
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
		//fmt.Printf(" %#v", f)
		listSockInner(f)
		return nil
	}
	return errorf("socket %d not found", inode)
}

// FIXME: an evil program could mess with this by dup()ing and close()ing a lot...
func ListAllFDs() error {
	files, err := ioutil.ReadDir(fdDir)
	if err != nil {
		return err
	}
	fmt.Println("---")
	defer fmt.Println("---")
	for _, file := range files {
		bname := path.Join(fdDir, file.Name())
		name, err := os.Readlink(bname)
		if err != nil {
			if err.(*os.PathError).Err == syscall.ENOENT {
				fmt.Printf("%s disappeared\n", file.Name())
			} else {
				return err
			}
			continue
		}
		fmt.Printf("%s -> %s", file.Name(), name)
		i, err := strconv.Atoi(file.Name())
		if err != nil {
			panic(err)
		}
		// FIXME: filename could have this format, but obvs would then not belong to the scheme...
		scheme := strings.SplitN(name, ":", 2)
		if len(scheme) > 1 {
			f := handlers[scheme[0]]
			if f != nil {
				i, err := os.Lstat(bname)
				if err != nil {
					fmt.Printf(" lstat failed: %s", err)
				} else {
					stat := i.Sys().(*syscall.Stat_t)
					f(stat, scheme[1])
				}
			} else {
				fmt.Printf(" no handler for '%s'", scheme[0])
			}
		}
		PrintRights(i)
		fmt.Println("")
	}
	return nil
}
