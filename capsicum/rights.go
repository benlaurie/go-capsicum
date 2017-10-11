package capsicum

import (
	"bytes"
	"fmt"
)

type capDesc struct {
	cap  uint64
	desc string
}

func (r *CapRights) String() string {
	if r == nil {
		return "[no rights]"
	}
	var b bytes.Buffer
	for _, d := range capDescs {
		s, err := CapRightsIsSet(r, d.cap)
		if err != nil {
			panic(err)
		}
		if s {
			b.WriteString(" ")
			b.WriteString(d.desc)
		}
	}
	return b.String()[1:]
}

func PrintRights(fd int) {
	r, err := CapRightsGetFd(uintptr(fd))
	if err != nil {
		panic(err)
	}
	fmt.Print(r.String())
}
