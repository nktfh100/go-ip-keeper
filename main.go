package main

import (
	_ "embed"
	"github.com/brotherpowers/ipsubnet"
	"github.com/rodrigocfd/windigo/win"
	"github.com/rodrigocfd/windigo/win/co"
	"regexp"
	"runtime"
)

type IPKeeper struct {
	mw     *MainWindow
	data   DataIPs
	subnet *ipsubnet.Ip
}

var IPRegex = regexp.MustCompile(`^((25[0-5]|(2[0-4]|1\d|[1-9]|)\d)\.?\b){4}$`)

//go:embed meme.jpg
var meme []byte

func main() {
	runtime.LockOSThread()

	programData := &IPKeeper{}

	programData.data = LoadData()

	if programData.data.IP != "" {
		programData.updateSubnet(programData.data.IP)
	}

	programData.mw = NewMainWindow(programData)

	programData.mw.wnd.RunAsMain()
}

func (programData *IPKeeper) updateSubnet(ip string) {
	programData.subnet = ipsubnet.SubnetCalculator(ip, 24)
}

func ShowError(hwnd win.HWND, text string) {
	hwnd.MessageBox(text, "ERROR", co.MB_ICONERROR)
}
