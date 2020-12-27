## TODO: need to build a good makefile

.PHONY: test

test:
	go build && ./signal-exporter export