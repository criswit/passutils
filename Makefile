.PHONY: exporter

exporter:
	mkdir -p bin
	go build -o bin/exporter ./exporter