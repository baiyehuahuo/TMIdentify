package main

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"identifyTencentMeeting/config"
	"identifyTencentMeeting/consts"
	"identifyTencentMeeting/types"
	"log"
	"os"
	"os/signal"
)

func main() {
	SetupCloseHandler()
	Initialize()
	defer Destroy()

	handle, err := pcap.OpenLive(config.GlobalConfig.DeviceName, consts.SnapLen, consts.Promiscuous, consts.TimeOut)
	fmt.Printf("I open the device: %v\n", config.GlobalConfig.DeviceName)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		types.GetFlowFromPacket(packet)
	}
}

// SetupCloseHandler 测试用，Ctrl+C打断程序时 打印cache中的消息
func SetupCloseHandler() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		Destroy()
		os.Exit(0)
	}()
}
