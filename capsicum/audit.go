package capsicum

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
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

type handler map[string]func(*syscall.Stat_t)

var handlers = handler{
	"socket": listSock,
}

func listSock(s *syscall.Stat_t) {
	fmt.Printf(" inode %d", s.Ino)
	f, err := os.Open(tcp6)
	if err != nil {
		log.Panic(err)
	}
	r := bufio.NewReader(f)
	l, p, err := r.ReadLine()
	if err != nil {
		panic(err)
	}
	if p {
		log.Panic("line too long")
	}
	if string(l) != "  sl  local_address                         remote_address                        st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode" {
		log.Panic("unknown format: %s", l)
	}
	for l, p, err = r.ReadLine(); ; l, p, err = r.ReadLine() {
		if err != nil {
			log.Panic(err)
		}
		c := strings.Fields(string(l))
		if len(c) < 11 {
			log.Panic("Don't understand: %s", l)
		}
		fmt.Printf(" %#v", c)
	}
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
					f(stat)
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
