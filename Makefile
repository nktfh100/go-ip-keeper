build:
	go build -ldflags="-H=windowsgui -s" -o "./bin/ip-keeper.exe"
	upx -9 -k ./bin/ip-keeper.exe