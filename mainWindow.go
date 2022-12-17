package main

import (
	"github.com/rodrigocfd/windigo/ui"
	"github.com/rodrigocfd/windigo/ui/wm"
	"github.com/rodrigocfd/windigo/win"
	"github.com/rodrigocfd/windigo/win/co"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

type MainWindow struct {
	wnd             ui.WindowMain
	startingIpInput ui.Edit
	btnNew          ui.Button
	btnEdit         ui.Button
	btnDel          ui.Button
	ipsList         ui.ListView
	btnsActive      bool
	memeBtn         ui.Button
}

func NewMainWindow(programData *IPKeeper) *MainWindow {
	wnd := ui.NewWindowMain(
		ui.WindowMainOpts().
			Title("IP keeper by nktfh100").
			ClientArea(win.SIZE{Cx: 360, Cy: 380}).
			IconId(101),
	)

	me := &MainWindow{
		wnd: wnd,
		ipsList: ui.NewListView(
			wnd,
			ui.ListViewOpts().
				Position(win.POINT{X: 10, Y: 10}).
				Size(win.SIZE{Cx: 230, Cy: 360}).
				CtrlExStyles(co.LVS_EX_FULLROWSELECT|co.LVS_EX_AUTOSIZECOLUMNS|co.LVS_EX_ONECLICKACTIVATE).
				CtrlStyles(co.LVS_REPORT|co.LVS_SHOWSELALWAYS|co.LVS_SINGLESEL),
		),
		btnNew: ui.NewButton(
			wnd,
			ui.ButtonOpts().
				Position(win.POINT{X: 250, Y: 20}).
				Size(win.SIZE{Cx: 100, Cy: 30}).
				Text("New"),
		),
		btnEdit: ui.NewButton(
			wnd,
			ui.ButtonOpts().
				Position(win.POINT{X: 250, Y: 60}).
				Size(win.SIZE{Cx: 100, Cy: 30}).
				Text("Edit"),
		),
		btnDel: ui.NewButton(
			wnd,
			ui.ButtonOpts().
				Position(win.POINT{X: 250, Y: 100}).
				Size(win.SIZE{Cx: 100, Cy: 30}).
				Text("Delete"),
		),
		// Add "No DHCP?" button at the bottom
		memeBtn: ui.NewButton(
			wnd,
			ui.ButtonOpts().
				Position(win.POINT{X: 250, Y: 320}).
				Size(win.SIZE{Cx: 100, Cy: 30}).
				Text("By nktfh100"),
		),
	}

	me.ipsList.On().LvnItemChanged(func(p *win.NMLISTVIEW) {
		if p.UNewState == co.LVIS_NONE && me.btnsActive {
			// None selected
			me.btnEdit.Hwnd().EnableWindow(false)
			me.btnDel.Hwnd().EnableWindow(false)
			me.btnsActive = false
		} else if !me.btnsActive {
			me.btnEdit.Hwnd().EnableWindow(true)
			me.btnDel.Hwnd().EnableWindow(true)
			me.btnsActive = true
		}
	})

	wnd.On().WmCreate(func(p wm.Create) int {

		// Disable the buttons until an item is selected
		me.btnEdit.Hwnd().EnableWindow(false)
		me.btnDel.Hwnd().EnableWindow(false)

		// Add the list columns
		me.ipsList.Columns().Add([]int{100, 120}, "IP", "Owner")

		// Add all the ips to the list
		for _, ip := range programData.data.IPs {
			me.ipsList.Items().Add(ip.IP, ip.Owner)
		}

		if programData.data.IP == "" {
			go func() {
				OpenIPModal(programData)
			}()
		}

		return 0
	})

	me.btnNew.On().BnClicked(func() {
		OpenAddIPModal(programData)
	})

	me.btnEdit.On().BnClicked(func() {
		selectedItems := me.ipsList.Items().SelectedItems()
		if len(selectedItems) != 1 {
			ShowError(me.wnd.Hwnd(), "No IP selected!")
			return
		}
		selectedItem := selectedItems[0]
		selectedItemIP := selectedItem.Text(0)
		var selectedItemI int
		// Get the IP index
		for i, ele := range programData.data.IPs {
			if ele.IP == selectedItemIP {
				selectedItemI = i
				break
			}
		}
		OpenEditIPWindow(programData, selectedItemI, selectedItem)
	})

	me.btnDel.On().BnClicked(func() {
		selectedItems := me.ipsList.Items().SelectedItems()
		if len(selectedItems) != 1 {
			ShowError(me.wnd.Hwnd(), "No IP selected!")
			return
		}
		selectedItem := selectedItems[0]
		selectedItemIP := selectedItem.Text(0)
		var selectedItemI int
		// Get the IP index
		for i, ele := range programData.data.IPs {
			if ele.IP == selectedItemIP {
				selectedItemI = i
				break
			}
		}
		// Remove the IP based on the index
		programData.data.removeIPByI(selectedItemI)
		selectedItem.Delete()

		SaveData(programData.data)

		me.btnEdit.Hwnd().EnableWindow(false)
		me.btnDel.Hwnd().EnableWindow(false)
		me.btnsActive = false
	})

	me.memeBtn.On().BnClicked(func() {
		me.memeBtn.Hwnd().ShowWindow(co.SW_HIDE)
		// Write the meme to temp folder
		memePath := filepath.Join(os.TempDir(), "meme.png")
		err := ioutil.WriteFile(memePath, meme, 0644)
		if err != nil {
			return
		}

		exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", memePath).Start()
	})

	return me
}
