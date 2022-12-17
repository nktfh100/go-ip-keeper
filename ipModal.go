package main

import (
	"fmt"
	"github.com/rodrigocfd/windigo/ui"
	"github.com/rodrigocfd/windigo/ui/wm"
	"github.com/rodrigocfd/windigo/win"
	"github.com/rodrigocfd/windigo/win/co"
	"os"
)

func OpenIPModal(programData *IPKeeper) {
	wnd := ui.NewWindowModal(
		ui.WindowModalOpts().
			Title("IP address").
			ClientArea(win.SIZE{Cx: 200, Cy: 110}),
	)

	ipInput := ui.NewEdit(
		wnd,
		ui.EditOpts().
			Position(win.POINT{X: 10, Y: 40}).
			Size(win.SIZE{Cx: 180, Cy: 20}),
	)

	submitBtn := ui.NewButton(
		wnd,
		ui.ButtonOpts().
			Position(win.POINT{X: 180/2 + 10 - 50, Y: 70}).
			Size(win.SIZE{Cx: 100, Cy: 30}).
			Text("Submit"),
	)

	ui.NewStatic(
		wnd,
		ui.StaticOpts().
			Position(win.POINT{X: 10, Y: 10}).
			Size(win.SIZE{Cx: 180, Cy: 20}).
			Text("Enter the starting IP address:").
			CtrlStyles(co.SS_CENTER),
	)

	wnd.On().WmCreate(func(p wm.Create) int {
		wnd.Hwnd().SetWindowLongPtr(co.GWLP_STYLE, uintptr(co.WS_CAPTION)|uintptr(co.WS_SYSMENU))
		return 0
	})

	wnd.On().WmClose(func() {
		if programData.data.IP == "" {
			os.Exit(0)
		} else {
			wnd.Hwnd().DestroyWindow()
			programData.mw.wnd.Hwnd().EnableWindow(true)
		}
	})

	submitBtn.On().BnClicked(func() {
		IPText := ipInput.Text()
		// Check if the IP provided is a valid IP address
		isOk := IPRegex.MatchString(IPText)
		if !isOk {
			msg := fmt.Sprintf("%s is not a valid IP!", IPText)
			ShowError(wnd.Hwnd(), msg)
			return
		}
		programData.data.IP = IPText
		programData.updateSubnet(programData.data.IP)
		SaveData(programData.data)
		wnd.Hwnd().PostMessage(co.WM_CLOSE, 0, 0)
	})

	// This for some reason will stay blocked even after the window is closed
	wnd.ShowModal(programData.mw.wnd)
}
