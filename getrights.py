#!/usr/bin/env python

import string
import sys

symbols = {}
with open(sys.argv[1]) as f:
    for line in f:
        #print line
        if line.startswith("#define CAP_"):
            s = line.split()[1]
            if not s.endswith(")") and not s.startswith("CAP_RIGHTS_VERSION") and not s.startswith("CAP_FCNTL_") and not s.endswith("_ALL"):
                symbols[s] = True

ss = sorted(symbols)

print """package capsicum

// #include <sys/capsicum.h>
import "C"
"""

print "const ("

for s in ss:
    print "\t%s = C.%s" % (s, s)
print ")"

print "var capDescs = [...]capDesc{"
for s in ss:
    print '\t{%s, "%s"},' % (s, s[4:].lower())
print "}"

