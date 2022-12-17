package main

import (
	"fmt"
	"github.com/rodrigocfd/windigo/ui"
	"github.com/rodrigocfd/windigo/ui/wm"
	"github.com/rodrigocfd/windigo/win"
	"github.com/rodrigocfd/windigo/win/co"
)

func OpenEditIPWindow(programData *IPKeeper, ipI int, selectedItem ui.ListViewItem) {
	wnd := ui.NewWindowModal(
		ui.WindowModalOpts().
			Title("Edit IP").
			ClientArea(win.SIZE{Cx: 160, Cy: 150}),
	)

	ipInput := ui.NewEdit(
		wnd,
		ui.EditOpts().
			Position(win.POINT{X: 10, Y: 35}).
			Size(win.SIZE{Cx: 140, Cy: 20}).
			Text(selectedItem.Text(0)),
	)
	ownerInput := ui.NewEdit(
		wnd,
		ui.EditOpts().
			Position(win.POINT{X: 10, Y: 80}).
			Size(win.SIZE{Cx: 140, Cy: 20}).
			Text(selectedItem.Text(1)),
	)
	submitBtn := ui.NewButton(
		wnd,
		ui.ButtonOpts().
			Position(win.POINT{X: 140/2 + 10 - 50, Y: 110}).
			Size(win.SIZE{Cx: 100, Cy: 30}).
			Text("Confirm"),
	)

	ui.NewStatic(
		wnd,
		ui.StaticOpts().
			Position(win.POINT{X: 10, Y: 15}).
			Size(win.SIZE{Cx: 140, Cy: 20}).
			Text("IP:"),
	)

	ui.NewStatic(
		wnd,
		ui.StaticOpts().
			Position(win.POINT{X: 10, Y: 60}).
			Size(win.SIZE{Cx: 140, Cy: 20}).
			Text("Owner:"),
	)

	wnd.On().WmCreate(func(p wm.Create) int {
		wnd.Hwnd().SetWindowLongPtr(co.GWLP_STYLE, uintptr(co.WS_CAPTION)|uintptr(co.WS_SYSMENU))
		return 0
	})

	submitBtn.On().BnClicked(func() {
		IPText := ipInput.Text()
		ownerText := ownerInput.Text()
		// Check if the IP provided is a valid IP address
		if !IPRegex.MatchString(IPText) {
			msg := fmt.Sprintf("%s is not a valid IP!", IPText)
			ShowError(wnd.Hwnd(), msg)
			return
		}

		if ownerText == "" {
			ShowError(wnd.Hwnd(), "Owner field cannot be empty!")
			return
		}

		if IPText != selectedItem.Text(0) && programData.data.ipExists(IPText) {
			ShowError(wnd.Hwnd(), "IP already exists!")
			return
		}

		programData.data.IPs[ipI].IP = IPText
		programData.data.IPs[ipI].Owner = ownerText

		selectedItem.SetText(0, IPText)
		selectedItem.SetText(1, ownerText)

		SaveData(programData.data)
		wnd.Hwnd().PostMessage(co.WM_CLOSE, 0, 0)
	})

	// This for some reason will stay blocked even after the window is closed
	wnd.ShowModal(programData.mw.wnd)
}
