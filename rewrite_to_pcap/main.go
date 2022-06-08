// 复写当前目录下所有 pcap 文件
package main

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/pcapgo"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	SnapShotLen uint32 = 1024
	TempPcap           = "temp.pcap"

	MLNum  = 5
	DPINum = 3
	Saved  = 50
)

var (
	count = 0
	total = 0
)

type Flow struct {
	packets    []gopacket.Packet
	dataPacket int
}

func main() {
	// dirPath := "./one/two"
	// fmt.Println()
	nums := "23"
	b := byte('b')
	a := nums[0] - '0'
	fmt.Printf("%v\t%v\t%v", '0', b, a-b)
	var err error
	if err = filepath.Walk("NonVPN_50", func(filePath string, info os.FileInfo, err error) error {
		if path.Ext(filePath) == ".pcap" {
			if err = truncatePure(filePath); err != nil {
				fmt.Printf("file %s rewrite fail: %v\n", filePath, err)
			}
		}
		return nil
	}); err != nil {
		log.Fatal(err)
	}
	fmt.Println(count, total)
}

func reWrite(pcapFilePath string) (err error) {
	var handle *pcap.Handle
	handle, err = pcap.OpenOffline(pcapFilePath)
	pcapFilePath = changePath(pcapFilePath)
	if err != nil {
		// fmt.Println(err)
		return nil
	}

	f, _ := os.Create(TempPcap)
	w := pcapgo.NewWriter(f)
	if err = w.WriteFileHeader(SnapShotLen, layers.LinkTypeEthernet); err != nil {
		return err
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		if err = w.WritePacket(packet.Metadata().CaptureInfo, packet.Data()); err != nil {
			return err
		}
	}
	handle.Close()
	f.Close()
	// if err = os.Remove(pcapFilePath); err != nil {
	// 	log.Println(err)
	// }
	ext := path.Ext(pcapFilePath)
	if err = os.Rename(TempPcap, pcapFilePath[:len(pcapFilePath)-len(ext)]+".pcap"); err != nil {
		_ = os.Remove(TempPcap)
		log.Println(err)
	}
	return nil
}

func truncate(pcapFilePath string) (err error) {
	flowMap := map[string]int{}
	flowDataMap := map[string]int{}
	var handle *pcap.Handle
	handle, err = pcap.OpenOffline(pcapFilePath)
	pcapFilePath = changePath(pcapFilePath)
	if err != nil {
		// fmt.Println(err)
		return nil
	}

	f, _ := os.Create(TempPcap)
	w := pcapgo.NewWriter(f)
	if err = w.WriteFileHeader(SnapShotLen, layers.LinkTypeEthernet); err != nil {
		return err
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		flowKey := getFlowKey(packet)
		flowMap[flowKey]++
		if app := packet.ApplicationLayer(); app != nil {
			flowDataMap[flowKey]++
		}
		if flowKey == "" || (flowMap[flowKey] >= Saved && flowDataMap[flowKey] >= DPINum) {
			continue
		}
		if err = w.WritePacket(packet.Metadata().CaptureInfo, packet.Data()); err != nil {
			return err
		}
	}
	handle.Close()
	f.Close()
	ext := path.Ext(pcapFilePath)
	if err = os.Rename(TempPcap, pcapFilePath[:len(pcapFilePath)-len(ext)]+"_"+strconv.Itoa(Saved)+".pcap"); err != nil {
		_ = os.Remove(TempPcap)
		log.Println(err)
	}
	return nil
}

func truncatePure(pcapFilePath string) (err error) {
	flowMap := map[string]*Flow{}
	var handle *pcap.Handle
	handle, err = pcap.OpenOffline(pcapFilePath)
	pcapFilePath = changePath(pcapFilePath)
	if err != nil {
		// fmt.Println(err)
		return nil
	}

	f, _ := os.Create(TempPcap)
	w := pcapgo.NewWriter(f)
	if err = w.WriteFileHeader(SnapShotLen, layers.LinkTypeEthernet); err != nil {
		return err
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	var flowKey string
	for packet := range packetSource.Packets() {
		if flowKey = getFlowKey(packet); flowKey == "" {
			continue
		}
		if flowMap[flowKey] == nil {
			flowMap[flowKey] = &Flow{}
		}
		flowMap[flowKey].addPacket(packet, Saved)
	}
	total += len(flowMap)
	for _, flow := range flowMap {
		if flow.dataPacket >= DPINum && len(flow.packets) >= MLNum {
			for _, packet := range flow.packets {
				_ = w.WritePacket(packet.Metadata().CaptureInfo, packet.Data())
			}
			count++
		}
	}
	handle.Close()
	f.Close()
	// if err = os.Remove(pcapFilePath); err != nil {
	// 	log.Println(err)
	// }
	ext := path.Ext(pcapFilePath)
	if err = os.Rename(TempPcap, pcapFilePath[:len(pcapFilePath)-len(ext)]+"_pure"+".pcap"); err != nil {
		_ = os.Remove(TempPcap)
		log.Println(err)
	}
	return nil
}

func (flow *Flow) addPacket(packet gopacket.Packet, maxSaved int) {
	if len(flow.packets) >= maxSaved {
		return
	}
	flow.packets = append(flow.packets, packet)
	if app := packet.ApplicationLayer(); app != nil {
		flow.dataPacket++
	}
}

// getFlowKey get a unique key from packet by ip:port
func getFlowKey(packet gopacket.Packet) string {
	networkLayer := packet.NetworkLayer()
	transportLayer := packet.TransportLayer()
	if networkLayer == nil || transportLayer == nil {
		return ""
	}
	network, transport := networkLayer.NetworkFlow(), transportLayer.TransportFlow()
	transportLayerType := transportLayer.LayerType()

	srcIP, dstIP := network.Endpoints()
	srcPort, dstPort := transport.Endpoints()

	if srcIP.LessThan(dstIP) {
		return fmt.Sprintf("%s: %s:%s,%s:%s", transportLayerType, srcIP, srcPort, dstIP, dstPort)
	}
	return fmt.Sprintf("%s: %s:%s,%s:%s", transportLayerType, dstIP, dstPort, srcIP, srcPort)
}

func changePath(pathStr string) string {
	paths := strings.Split(pathStr, "\\")
	paths[0] += "_" + "pure" // strconv.Itoa(Saved)
	result := strings.Join(paths, "\\")
	// fmt.Println(path.Dir(result))
	dirPath := strings.Join(paths[:len(paths)-1], "\\")
	if _, err := os.Stat(dirPath); err != nil {
		if err = os.MkdirAll(dirPath, 0777); err != nil {
			log.Fatal(err)
		}
	}
	return result
}
