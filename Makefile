xiaocao.exe: main.go
	go build -ldflags="-H windowsgui -w -s" -o xiaocao.exe

syso: 1024.manifest icon.ico
	rsrc -manifest 1024.manifest -ico icon.ico -o 1024.syso

.PHONY: clean gverInfo

gverInfo: icon.ico 1024.manifest
	go get github.com/josephspurrier/goversioninfo/cmd/goversioninfo
	goversioninfo -icon=icon.ico -manifest=1024.manifest

clean: xiaocao.exe
	-rm xiaocao.exe