.PHONY: dev
dev:
	go build -o ./.tmp/main . && air dev
