# go-capsicum

Capsicum shim and utilities for Go.

Although Go includes support for Capsicum, it is currently only for
FreeBSD, and is missing some useful functions.

There is also a [version of Linux that supoorts
Capsicum](http://capsicum-linux.org/).

This library is intended to provide a universal shim for Capsicum and
also fill in missing functionality in the base system.
