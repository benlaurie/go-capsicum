// Cloned from https://golang.org/src/os/sticky_notbsd.go

// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !darwin
// +build !dragonfly
// +build !freebsd
// +build !netbsd
// +build !openbsd
// +build !solaris

package capsicum

const supportsCreateWithStickyBit = true
