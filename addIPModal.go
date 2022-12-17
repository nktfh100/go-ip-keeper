package main

import (
	"fmt"
	"github.com/rodrigocfd/windigo/ui"
	"github.com/rodrigocfd/windigo/ui/wm"
	"github.com/rodrigocfd/windigo/win"
	"github.com/rodrigocfd/windigo/win/co"
)

func OpenAddIPModal(programData *IPKeeper) {
	wnd := ui.NewWindowModal(
		ui.WindowModalOpts().
			Title("New IP").
			ClientArea(win.SIZE{Cx: 160, Cy: 150}),
	)

	ipInput := ui.NewEdit(
		wnd,
		ui.EditOpts().
			Position(win.POINT{X: 10, Y: 35}).
			Size(win.SIZE{Cx: 140, Cy: 20}),
	)

	ownerInput := ui.NewEdit(
		wnd,
		ui.EditOpts().
			Position(win.POINT{X: 10, Y: 80}).
			Size(win.SIZE{Cx: 140, Cy: 20}),
	)

	submitBtn := ui.NewButton(
		wnd,
		ui.ButtonOpts().
			Position(win.POINT{X: 140/2 + 10 - 50, Y: 110}).
			Size(win.SIZE{Cx: 100, Cy: 30}).
			Text("Add"),
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
		go func() {
			// Generate a new valid unused IP and set it to the ip field
			unusedIP := GetUnusedIP(programData)
			ipInput.SetText(unusedIP)
		}()
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

		if programData.data.ipExists(IPText) {
			ShowError(wnd.Hwnd(), "IP already exists!")
			return
		}

		newIPs := append(programData.data.IPs, DataIP{
			IP:    IPText,
			Owner: ownerText,
		})
		programData.data.IPs = newIPs
		programData.mw.ipsList.Items().Add(IPText, ownerText)

		SaveData(programData.data)
		wnd.Hwnd().PostMessage(co.WM_CLOSE, 0, 0)
	})

	// This for some reason will stay blocked even after the window is closed
	wnd.ShowModal(programData.mw.wnd)
}
