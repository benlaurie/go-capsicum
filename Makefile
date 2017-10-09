capsicum/rights_info.go: /usr/include/linux/capsicum.h getrights.py
	./getrights.py $< > $@
	go fmt $@
