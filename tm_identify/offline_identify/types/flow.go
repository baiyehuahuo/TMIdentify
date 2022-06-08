package types

import (
	"encoding/binary"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/patrickmn/go-cache"
	"identifyTencentMeeting/classify/DPI"
	"identifyTencentMeeting/classify/ML/decision_tree"
	"identifyTencentMeeting/config"
	"identifyTencentMeeting/consts"
	"identifyTencentMeeting/utils"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

type PacketFeature struct {
	PayloadLength int
	Timestamp     time.Time
}

type Flow struct {
	FlowStr           string
	PacketsFeature    []PacketFeature
	Payloads          []string
	ForwardPackets    []int             // location of forward packets in Flow.packets
	BackwardPackets   []int             // location of backward packets in Flow.packets
	ForwardIP         gopacket.Endpoint // The device ip used for initialization
	ForwardPort       gopacket.Endpoint
	BackwardIP        gopacket.Endpoint // Not the ip of the local device ip
	BackwardPort      gopacket.Endpoint
	TransportProtocol layers.IPProtocol
	HasRTCP           int
	IsTencentMeeting  bool
	mtx               sync.RWMutex
}

var (
	flowCache    *cache.Cache
	flowCacheMtx sync.Mutex
	T            int
	F            int
)

// AddPacket add a packet to a flow
func (flow *Flow) addPacket(packet gopacket.Packet) {
	flow.mtx.RLock()
	if !flow.IsTencentMeeting && len(flow.PacketsFeature) < config.GlobalConfig.MLJudgePacket {
		srcIP := packet.NetworkLayer().NetworkFlow().Src()
		srcPort := packet.TransportLayer().TransportFlow().Src()
		// fmt.Println(flow.srcIP, srcIP, flow.srcPort, srcPort)
		if srcIP == flow.ForwardIP && srcPort == flow.ForwardPort {
			flow.ForwardPackets = append(flow.ForwardPackets, len(flow.PacketsFeature))
		} else if srcIP == flow.BackwardIP && srcPort == flow.BackwardPort {
			flow.BackwardPackets = append(flow.BackwardPackets, len(flow.PacketsFeature))
		}
		packetFeature := PacketFeature{
			PayloadLength: 0,
			Timestamp:     packet.Metadata().Timestamp,
		}

		if app := packet.ApplicationLayer(); app != nil {
			packetFeature.PayloadLength = len(app.Payload())
		}
		flow.PacketsFeature = append(flow.PacketsFeature, packetFeature)
		if flow.HasRTCP == consts.WithoutRTCP {
			flow.HasRTCP = isRTCP(packet)
		}
		if len(flow.PacketsFeature) == config.GlobalConfig.MLJudgePacket {
			SaveFeaturesToCSV(flow.FlowStr, flow)
			if config.GlobalConfig.UseML {
				features := CalculateFeaturesFromFlow(flow)
				flow.IsTencentMeeting = flow.IsTencentMeeting || decision_tree.Classify(convertFeaturesToFloat64(features), "CART")
				if config.GlobalConfig.PrintDetect && flow.IsTencentMeeting {
					fmt.Printf("DT find: %s is tencent meeting flow\n", flow.FlowStr)
				}
				if config.GlobalConfig.SplitFlow {
					flowCache.Delete(flow.FlowStr)
				}
			}
		}
	}
	if app := packet.ApplicationLayer(); !flow.IsTencentMeeting && app != nil && len(flow.Payloads) < config.GlobalConfig.DPIJudgePacket {
		flow.Payloads = append(flow.Payloads, string(app.Payload()))
		if config.GlobalConfig.UseDPI && len(flow.Payloads) == config.GlobalConfig.DPIJudgePacket {
			// fmt.Printf("?")
			flow.IsTencentMeeting = flow.IsTencentMeeting || DPI.Classify(strings.Join(flow.Payloads, ""))
			if config.GlobalConfig.PrintDetect && flow.IsTencentMeeting {
				fmt.Printf("DPI find: %s is tencent meeting flow\n", flow.FlowStr)
			}
		}
	}
	flow.mtx.RUnlock()
}

// getFlowKey get a unique key from packet by ip:port
func getFlowKey(network, transport gopacket.Flow, layerType gopacket.LayerType) string {
	srcIP, dstIP := network.Endpoints()
	srcPort, dstPort := transport.Endpoints()
	if srcIP.LessThan(dstIP) {
		return fmt.Sprintf("%s: %s:%s,%s:%s", layerType, srcIP, srcPort, dstIP, dstPort)
	}
	return fmt.Sprintf("%s: %s:%s,%s:%s", layerType, dstIP, dstPort, srcIP, srcPort)
}

// NewFlow create a new flow from a packet
func newFlow(packet gopacket.Packet, flowStr string) *Flow {
	flow := &Flow{
		FlowStr:        flowStr,
		PacketsFeature: make([]PacketFeature, 0, config.GlobalConfig.MLJudgePacket),
		mtx:            sync.RWMutex{},
	}

	if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
		flow.TransportProtocol = layers.IPProtocolTCP
	} else if udpLayer := packet.Layer(layers.LayerTypeUDP); udpLayer != nil {
		flow.TransportProtocol = layers.IPProtocolUDP
	}

	flow.setForwardByFirstPacketSrc(packet)

	return flow
}

// setForwardByFirstPacketSrc
func (flow *Flow) setForwardByFirstPacketSrc(packet gopacket.Packet) {
	flow.ForwardIP, flow.BackwardIP = packet.NetworkLayer().NetworkFlow().Endpoints()
	flow.ForwardPort, flow.BackwardPort = packet.TransportLayer().TransportFlow().Endpoints()
}

// DestroyCache destroy the hash cache
func DestroyCache() {
	gConfig := config.GlobalConfig
	counter := func(flow *Flow) {
		if flow.IsTencentMeeting {
			T++
		} else {
			F++
		}
	}
	if flowCache != nil {
		// fmt.Println(len(flowCache.Items()))
		for key, item := range flowCache.Items() {
			flowCache.Delete(key)
			flow := item.Object.(*Flow)

			// if len(flow.PacketsFeature) == gConfig.MLJudgePacket && len(flow.Payloads) == gConfig.DPIJudgePacket {
			// 	counter((flow))
			// }
			if (gConfig.UseML && len(flow.PacketsFeature) == gConfig.MLJudgePacket) || (gConfig.UseDPI && len(flow.Payloads) == gConfig.DPIJudgePacket) {
				counter(flow)
			}
		}
		flowCache = nil
	}
}

// GetFlowFromPacket get a flow from a packet and judge isTencentMeeting or not
func GetFlowFromPacket(packet gopacket.Packet) (flow *Flow, isTencentMeeting bool) {
	network := packet.NetworkLayer()
	transport := packet.TransportLayer()
	if network != nil && transport != nil {
		flowStr := getFlowKey(network.NetworkFlow(), transport.TransportFlow(), transport.LayerType())
		flowCacheMtx.Lock()
		savedFlow, ok := flowCache.Get(flowStr)
		if ok {
			flow = savedFlow.(*Flow)
		} else {
			flow = newFlow(packet, flowStr)
			if flow == nil {
				flowCacheMtx.Unlock()
				return nil, false
			}
		}
		flowCache.Set(flowStr, flow, cache.DefaultExpiration)
		flowCacheMtx.Unlock()
		flow.addPacket(packet)
		isTencentMeeting = flow.IsTencentMeeting
	}
	return
}

// InitCache initialize a hash cache (map[string]*flow)
func InitCache(expiration time.Duration) {
	flowCache = cache.New(expiration, 5*time.Minute)
	// flowCache.OnEvicted(SaveFeaturesToCSV)
}

func ResetCache() (int, int) {
	DestroyCache()
	t, f := T, F
	T, F = 0, 0
	InitCache(consts.CacheExpiration)
	return t, f
}

// SaveFeaturesToCSV save the features calculated by the flow into the csv file
func SaveFeaturesToCSV(flowKey string, item interface{}) {
	var fieldsValue = make([]string, 0, config.GlobalConfig.CSVColumn)
	flow := item.(*Flow)
	if config.GlobalConfig.FilterShort && len(flow.PacketsFeature) < config.GlobalConfig.MLJudgePacket {
		return
	}
	fieldsValue = append(fieldsValue, flowKey)
	fieldsValue = append(fieldsValue, utils.GetStructValuesAsString(CalculateFeaturesFromFlow(flow))...)
	fieldsValue = append(fieldsValue, strconv.FormatBool(config.GlobalConfig.Label))
	if err := config.GlobalConfig.CSVWriter.Write(fieldsValue); err != nil {
		log.Println(utils.ErrWrapOrWithMessage(false, err))
	}
	config.GlobalConfig.CSVWriter.Flush()
}

func isRTCP(packet gopacket.Packet) int {
	var payload []byte
	if layer := packet.Layer(layers.LayerTypeUDP); layer != nil {
		payload = layer.LayerPayload()
	}
	if len(payload) == 0 {
		return consts.WithoutRTCP
	}
	//          0                   1                   2                   3
	//         0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	//        +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// header |V=2|P|    RC   |       PT      |             length            |
	//        +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	totalVersionIs2 := true
	totalPacketInRange := true
	totalLength := 0
	checkPacketTypeInRange := func(packetType byte) bool {
		return packetType == 0 || 192 <= packetType && packetType <= 195 || 200 <= packetType && packetType <= 213 || packetType == 255
	}
	for pos := 0; pos+4 < len(payload) && totalVersionIs2 && totalPacketInRange; {
		totalVersionIs2 = totalVersionIs2 && (payload[pos]>>6 == 2)
		totalPacketInRange = totalPacketInRange && checkPacketTypeInRange(payload[pos+1])
		totalLength += (int(binary.BigEndian.Uint16(payload[pos+2:pos+4])) + 1) * 4
		pos = totalLength
	}

	if totalVersionIs2 && totalPacketInRange && totalLength == len(payload) {
		return consts.WithRTCP
	}
	return consts.WithoutRTCP
}
