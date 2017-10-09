package capsicum

import "fmt"

type capDesc struct {
	cap  uint64
	desc string
}

func PrintRights(fd int) {
	for _, d := range capDescs {
		r, err := CapRightsGetFd(uintptr(fd))
		if err != nil {
			panic(err)
		}
		s, err := CapRightsIsSet(r, d.cap)
		if err != nil {
			panic(err)
		}
		if s {
			fmt.Print(" ", d.desc)
		}
	}
}
