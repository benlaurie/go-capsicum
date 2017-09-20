// Cloned from https://golang.org/src/os/sticky_bsd.go

// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd netbsd openbsd solaris

package capsicum

// According to sticky(8), neither open(2) nor mkdir(2) will create
// a file with the sticky bit set.
const supportsCreateWithStickyBit = false
