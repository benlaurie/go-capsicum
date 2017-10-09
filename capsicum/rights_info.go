package capsicum

// #include <sys/capsicum.h>
import "C"

const (
	CAP_ACCEPT           = C.CAP_ACCEPT
	CAP_ACL_CHECK        = C.CAP_ACL_CHECK
	CAP_ACL_DELETE       = C.CAP_ACL_DELETE
	CAP_ACL_GET          = C.CAP_ACL_GET
	CAP_ACL_SET          = C.CAP_ACL_SET
	CAP_ALL0             = C.CAP_ALL0
	CAP_ALL1             = C.CAP_ALL1
	CAP_BIND             = C.CAP_BIND
	CAP_BINDAT           = C.CAP_BINDAT
	CAP_BPF              = C.CAP_BPF
	CAP_CHFLAGSAT        = C.CAP_CHFLAGSAT
	CAP_CONNECT          = C.CAP_CONNECT
	CAP_CONNECTAT        = C.CAP_CONNECTAT
	CAP_CREATE           = C.CAP_CREATE
	CAP_EPOLL_CTL        = C.CAP_EPOLL_CTL
	CAP_EVENT            = C.CAP_EVENT
	CAP_EXTATTR_DELETE   = C.CAP_EXTATTR_DELETE
	CAP_EXTATTR_GET      = C.CAP_EXTATTR_GET
	CAP_EXTATTR_LIST     = C.CAP_EXTATTR_LIST
	CAP_EXTATTR_SET      = C.CAP_EXTATTR_SET
	CAP_FCHDIR           = C.CAP_FCHDIR
	CAP_FCHFLAGS         = C.CAP_FCHFLAGS
	CAP_FCHMOD           = C.CAP_FCHMOD
	CAP_FCHMODAT         = C.CAP_FCHMODAT
	CAP_FCHOWN           = C.CAP_FCHOWN
	CAP_FCHOWNAT         = C.CAP_FCHOWNAT
	CAP_FCNTL            = C.CAP_FCNTL
	CAP_FEXECVE          = C.CAP_FEXECVE
	CAP_FLOCK            = C.CAP_FLOCK
	CAP_FPATHCONF        = C.CAP_FPATHCONF
	CAP_FSCK             = C.CAP_FSCK
	CAP_FSIGNAL          = C.CAP_FSIGNAL
	CAP_FSTAT            = C.CAP_FSTAT
	CAP_FSTATAT          = C.CAP_FSTATAT
	CAP_FSTATFS          = C.CAP_FSTATFS
	CAP_FSYNC            = C.CAP_FSYNC
	CAP_FTRUNCATE        = C.CAP_FTRUNCATE
	CAP_FUTIMES          = C.CAP_FUTIMES
	CAP_FUTIMESAT        = C.CAP_FUTIMESAT
	CAP_GETPEERNAME      = C.CAP_GETPEERNAME
	CAP_GETSOCKNAME      = C.CAP_GETSOCKNAME
	CAP_GETSOCKOPT       = C.CAP_GETSOCKOPT
	CAP_IOCTL            = C.CAP_IOCTL
	CAP_KQUEUE           = C.CAP_KQUEUE
	CAP_KQUEUE_CHANGE    = C.CAP_KQUEUE_CHANGE
	CAP_KQUEUE_EVENT     = C.CAP_KQUEUE_EVENT
	CAP_LINKAT_SOURCE    = C.CAP_LINKAT_SOURCE
	CAP_LINKAT_TARGET    = C.CAP_LINKAT_TARGET
	CAP_LISTEN           = C.CAP_LISTEN
	CAP_LOOKUP           = C.CAP_LOOKUP
	CAP_MAC_GET          = C.CAP_MAC_GET
	CAP_MAC_SET          = C.CAP_MAC_SET
	CAP_MKDIRAT          = C.CAP_MKDIRAT
	CAP_MKFIFOAT         = C.CAP_MKFIFOAT
	CAP_MKNODAT          = C.CAP_MKNODAT
	CAP_MMAP             = C.CAP_MMAP
	CAP_MMAP_R           = C.CAP_MMAP_R
	CAP_MMAP_RW          = C.CAP_MMAP_RW
	CAP_MMAP_RWX         = C.CAP_MMAP_RWX
	CAP_MMAP_RX          = C.CAP_MMAP_RX
	CAP_MMAP_W           = C.CAP_MMAP_W
	CAP_MMAP_WX          = C.CAP_MMAP_WX
	CAP_MMAP_X           = C.CAP_MMAP_X
	CAP_NOTIFY           = C.CAP_NOTIFY
	CAP_PDGETPID         = C.CAP_PDGETPID
	CAP_PDGETPID_FREEBSD = C.CAP_PDGETPID_FREEBSD
	CAP_PDKILL           = C.CAP_PDKILL
	CAP_PDKILL_FREEBSD   = C.CAP_PDKILL_FREEBSD
	CAP_PDWAIT           = C.CAP_PDWAIT
	CAP_PEELOFF          = C.CAP_PEELOFF
	CAP_PERFMON          = C.CAP_PERFMON
	CAP_POLL_EVENT       = C.CAP_POLL_EVENT
	CAP_PREAD            = C.CAP_PREAD
	CAP_PWRITE           = C.CAP_PWRITE
	CAP_READ             = C.CAP_READ
	CAP_RECV             = C.CAP_RECV
	CAP_RENAMEAT_SOURCE  = C.CAP_RENAMEAT_SOURCE
	CAP_RENAMEAT_TARGET  = C.CAP_RENAMEAT_TARGET
	CAP_SEEK             = C.CAP_SEEK
	CAP_SEEK_TELL        = C.CAP_SEEK_TELL
	CAP_SEM_GETVALUE     = C.CAP_SEM_GETVALUE
	CAP_SEM_POST         = C.CAP_SEM_POST
	CAP_SEM_WAIT         = C.CAP_SEM_WAIT
	CAP_SEND             = C.CAP_SEND
	CAP_SETNS            = C.CAP_SETNS
	CAP_SETSOCKOPT       = C.CAP_SETSOCKOPT
	CAP_SHUTDOWN         = C.CAP_SHUTDOWN
	CAP_SOCK_CLIENT      = C.CAP_SOCK_CLIENT
	CAP_SOCK_SERVER      = C.CAP_SOCK_SERVER
	CAP_SYMLINKAT        = C.CAP_SYMLINKAT
	CAP_TTYHOOK          = C.CAP_TTYHOOK
	CAP_UNLINKAT         = C.CAP_UNLINKAT
	CAP_UNUSED0_44       = C.CAP_UNUSED0_44
	CAP_UNUSED0_57       = C.CAP_UNUSED0_57
	CAP_UNUSED1_27       = C.CAP_UNUSED1_27
	CAP_UNUSED1_57       = C.CAP_UNUSED1_57
	CAP_WRITE            = C.CAP_WRITE
)

var capDescs = [...]capDesc{
	{CAP_ACCEPT, "accept"},
	{CAP_ACL_CHECK, "acl_check"},
	{CAP_ACL_DELETE, "acl_delete"},
	{CAP_ACL_GET, "acl_get"},
	{CAP_ACL_SET, "acl_set"},
	{CAP_ALL0, "all0"},
	{CAP_ALL1, "all1"},
	{CAP_BIND, "bind"},
	{CAP_BINDAT, "bindat"},
	{CAP_BPF, "bpf"},
	{CAP_CHFLAGSAT, "chflagsat"},
	{CAP_CONNECT, "connect"},
	{CAP_CONNECTAT, "connectat"},
	{CAP_CREATE, "create"},
	{CAP_EPOLL_CTL, "epoll_ctl"},
	{CAP_EVENT, "event"},
	{CAP_EXTATTR_DELETE, "extattr_delete"},
	{CAP_EXTATTR_GET, "extattr_get"},
	{CAP_EXTATTR_LIST, "extattr_list"},
	{CAP_EXTATTR_SET, "extattr_set"},
	{CAP_FCHDIR, "fchdir"},
	{CAP_FCHFLAGS, "fchflags"},
	{CAP_FCHMOD, "fchmod"},
	{CAP_FCHMODAT, "fchmodat"},
	{CAP_FCHOWN, "fchown"},
	{CAP_FCHOWNAT, "fchownat"},
	{CAP_FCNTL, "fcntl"},
	{CAP_FEXECVE, "fexecve"},
	{CAP_FLOCK, "flock"},
	{CAP_FPATHCONF, "fpathconf"},
	{CAP_FSCK, "fsck"},
	{CAP_FSIGNAL, "fsignal"},
	{CAP_FSTAT, "fstat"},
	{CAP_FSTATAT, "fstatat"},
	{CAP_FSTATFS, "fstatfs"},
	{CAP_FSYNC, "fsync"},
	{CAP_FTRUNCATE, "ftruncate"},
	{CAP_FUTIMES, "futimes"},
	{CAP_FUTIMESAT, "futimesat"},
	{CAP_GETPEERNAME, "getpeername"},
	{CAP_GETSOCKNAME, "getsockname"},
	{CAP_GETSOCKOPT, "getsockopt"},
	{CAP_IOCTL, "ioctl"},
	{CAP_KQUEUE, "kqueue"},
	{CAP_KQUEUE_CHANGE, "kqueue_change"},
	{CAP_KQUEUE_EVENT, "kqueue_event"},
	{CAP_LINKAT_SOURCE, "linkat_source"},
	{CAP_LINKAT_TARGET, "linkat_target"},
	{CAP_LISTEN, "listen"},
	{CAP_LOOKUP, "lookup"},
	{CAP_MAC_GET, "mac_get"},
	{CAP_MAC_SET, "mac_set"},
	{CAP_MKDIRAT, "mkdirat"},
	{CAP_MKFIFOAT, "mkfifoat"},
	{CAP_MKNODAT, "mknodat"},
	{CAP_MMAP, "mmap"},
	{CAP_MMAP_R, "mmap_r"},
	{CAP_MMAP_RW, "mmap_rw"},
	{CAP_MMAP_RWX, "mmap_rwx"},
	{CAP_MMAP_RX, "mmap_rx"},
	{CAP_MMAP_W, "mmap_w"},
	{CAP_MMAP_WX, "mmap_wx"},
	{CAP_MMAP_X, "mmap_x"},
	{CAP_NOTIFY, "notify"},
	{CAP_PDGETPID, "pdgetpid"},
	{CAP_PDGETPID_FREEBSD, "pdgetpid_freebsd"},
	{CAP_PDKILL, "pdkill"},
	{CAP_PDKILL_FREEBSD, "pdkill_freebsd"},
	{CAP_PDWAIT, "pdwait"},
	{CAP_PEELOFF, "peeloff"},
	{CAP_PERFMON, "perfmon"},
	{CAP_POLL_EVENT, "poll_event"},
	{CAP_PREAD, "pread"},
	{CAP_PWRITE, "pwrite"},
	{CAP_READ, "read"},
	{CAP_RECV, "recv"},
	{CAP_RENAMEAT_SOURCE, "renameat_source"},
	{CAP_RENAMEAT_TARGET, "renameat_target"},
	{CAP_SEEK, "seek"},
	{CAP_SEEK_TELL, "seek_tell"},
	{CAP_SEM_GETVALUE, "sem_getvalue"},
	{CAP_SEM_POST, "sem_post"},
	{CAP_SEM_WAIT, "sem_wait"},
	{CAP_SEND, "send"},
	{CAP_SETNS, "setns"},
	{CAP_SETSOCKOPT, "setsockopt"},
	{CAP_SHUTDOWN, "shutdown"},
	{CAP_SOCK_CLIENT, "sock_client"},
	{CAP_SOCK_SERVER, "sock_server"},
	{CAP_SYMLINKAT, "symlinkat"},
	{CAP_TTYHOOK, "ttyhook"},
	{CAP_UNLINKAT, "unlinkat"},
	{CAP_UNUSED0_44, "unused0_44"},
	{CAP_UNUSED0_57, "unused0_57"},
	{CAP_UNUSED1_27, "unused1_27"},
	{CAP_UNUSED1_57, "unused1_57"},
	{CAP_WRITE, "write"},
}
