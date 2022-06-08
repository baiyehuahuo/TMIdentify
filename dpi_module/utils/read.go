package utils

import (
	"dpi_module/types"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"io/fs"
	"log"
	"path"
	"path/filepath"
)

func ReadPcapFile(pcapPath string) {
	handle, err := pcap.OpenOffline(pcapPath)
	if err != nil {
		log.Fatal(err)
	}
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		types.ClassifyPacket(packet)
	}
}

func ReadPcapDir(dirPath string) {
	var err error
	// fmt.Println(dirPath)
	if flowMap, flowPayloads := LoadFlows(dirPath); flowMap != nil && flowPayloads != nil {
		types.SetFlow(flowMap, flowPayloads)
	} else if err = filepath.Walk(dirPath, func(filePath string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			fmt.Println("dir is reading.", filePath)
		}
		if path.Ext(filePath) == ".pcap" {
			// fmt.Printf("reading pcap: %v\n", filePath)
			ReadPcapFile(filePath)
			types.ResetMap()
		}
		return nil
	}); err != nil {
		log.Fatal(err)
	} else {
		flowMap, flowPayloads = types.GetFlowMap(), types.GetFlowPayloads()
		SaveFlows(dirPath, flowMap, flowPayloads)
	}

}
