package main

import (
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

func GetUnusedIP(programData *IPKeeper) string {
	rangeIPs := programData.subnet.GetIPAddressRange()

	startIPS := strings.Split(rangeIPs[0], ".")
	endIPS := strings.Split(rangeIPs[1], ".")

	// Remove the last number from the IP
	baseIP := strings.Join(startIPS[:len(startIPS)-1], ".")

	firstNum, err := strconv.Atoi(startIPS[len(startIPS)-1])
	lastNum, err1 := strconv.Atoi(endIPS[len(endIPS)-1])

	if err != nil || err1 != nil {
		ShowError(programData.mw.wnd.Hwnd(), "Error converting IP number to int")
		return ""
	}

	var unusedIP string

	arpCmdData := executeArpCmd(programData)

	for i := firstNum; i <= lastNum; i++ {
		if i == 0 || i == 255 {
			continue
		}
		ip := baseIP + "." + strconv.Itoa(i)

		if !programData.data.ipExists(ip) {
			if strings.Contains(arpCmdData, ip) {
				continue
			} else {
				unusedIP = ip
				break
			}
		}
	}

	return unusedIP
}

func executeArpCmd(programData *IPKeeper) string {
	cmd := exec.Command("arp", "-a")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true} // To hide the console window
	stdout, err := cmd.Output()

	if err != nil {
		ShowError(programData.mw.wnd.Hwnd(), "Error executing arp command")
		return ""
	}

	return string(stdout)
}

//
//func PingIP(ip string) {
//	pinger, err := ping.NewPinger(ip)
//	pinger.SetPrivileged(true)
//	if err != nil {
//		panic(err)
//	}
//
//	pinger.Count = 2
//
//	//c := make(chan os.Signal, 1)
//	//signal.Notify(c, os.Interrupt)
//	//go func() {
//	//	for _ = range c {
//	//		pinger.Stop()
//	//	}
//	//}()
//
//	pinger.OnRecv = func(pkt *ping.Packet) {
//		fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v\n",
//			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
//	}
//	pinger.OnFinish = func(stats *ping.Statistics) {
//		fmt.Printf("\n--- %s ping statistics ---\n", stats.Addr)
//		fmt.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n",
//			stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
//		fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
//			stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
//	}
//
//	fmt.Printf("PING %s (%s):\n", pinger.Addr(), pinger.IPAddr())
//	err = pinger.Run()
//
//	if err != nil {
//		fmt.Printf("Failed to ping: %s\n", err.Error())
//	}
//}
