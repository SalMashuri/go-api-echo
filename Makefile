go-api-osx: main.go
	GOOS=darwin GOARCH=amd64 go build -ldflags '-s -w' -o $@