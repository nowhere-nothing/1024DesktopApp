package main

//go:generate rsrc -manifest 1024.manifest -ico icon.ico -o 1024.syso

/// go:generate go build -ldflags="-H windowsgui -w -s" -o xiaocao.exe

/// go get github.com/josephspurrier/goversioninfo/cmd/goversioninfo
/// go:generate goversioninfo -icon=icon.ico -manifest=1024.manifest
