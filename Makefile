.DEFAULT_GOAL := run

run:
	go run *.go

install:
	go install

stats:
	gocloc .


