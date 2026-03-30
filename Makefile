.PHONY: dev
dev:
	go build -o ./.tmp/main . && air dev localhost 8082
