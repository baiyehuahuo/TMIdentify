package types

import (
	"bytes"
	"dpi_module/config"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type Flow struct {
	FlowStr           string
	Payloads          []string
	FlowPayloads      string
	ForwardIP         string // The device ip used for initialization
	ForwardPort       string
	BackwardIP        string // Not the ip of the local device ip
	BackwardPort      string
	TransportProtocol layers.IPProtocol
}

var flowMap = map[string]*Flow{}
var flowPayloads []string
var total int

// AddPacket add a packet to a flow
func (flow *Flow) addPacket(packet gopacket.Packet) {
	if len(flow.Payloads) > config.GlobalConfig.JudgeThreshold {
		return
	}
	// fmt.Println(flow.srcIP, srcIP, flow.srcPort, srcPort)

	if app := packet.ApplicationLayer(); app != nil {
		flow.Payloads = append(flow.Payloads, string(app.Payload()))
		if len(flow.Payloads) == config.GlobalConfig.JudgeThreshold {
			buffer := bytes.Buffer{}
			for _, payload := range flow.Payloads[:config.GlobalConfig.JudgeThreshold] {
				buffer.WriteString(payload)
			}
			str := buffer.String()
			flow.FlowPayloads = str
			flowPayloads = append(flowPayloads, str)
		}
	}

}

// getFlowKey get a unique key from packet by ip:port
func getFlowKey(network, transport gopacket.Flow, transportLayerType gopacket.LayerType) string {
	srcIP, dstIP := network.Endpoints()
	srcPort, dstPort := transport.Endpoints()

	if srcIP.LessThan(dstIP) {
		return fmt.Sprintf("%s: %s:%s,%s:%s", transportLayerType, srcIP, srcPort, dstIP, dstPort)
	}
	return fmt.Sprintf("%s: %s:%s,%s:%s", transportLayerType, dstIP, dstPort, srcIP, srcPort)
}

// NewFlow create a new flow from a packet
func newFlow(packet gopacket.Packet, flowStr string) *Flow {
	flow := &Flow{
		FlowStr:  flowStr,
		Payloads: make([]string, 0, config.GlobalConfig.JudgeThreshold),
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
	forwardIP, backwardIP := packet.NetworkLayer().NetworkFlow().Endpoints()
	flow.ForwardIP, flow.BackwardIP = forwardIP.String(), backwardIP.String()
	forwardPort, backwardPort := packet.TransportLayer().TransportFlow().Endpoints()
	flow.ForwardPort, flow.BackwardPort = forwardPort.String(), backwardPort.String()
}

// ClassifyPacket Classify packets into a stream
func ClassifyPacket(packet gopacket.Packet) {
	network := packet.NetworkLayer()
	transport := packet.TransportLayer()
	if network != nil && transport != nil {
		flowStr := getFlowKey(network.NetworkFlow(), transport.TransportFlow(), transport.LayerType())
		flow, ok := flowMap[flowStr]
		if !ok {
			flow = newFlow(packet, flowStr)
		}
		flowMap[flowStr] = flow
		flow.addPacket(packet)
	}
}

func (flow *Flow) String() string {
	return fmt.Sprintf("%s:%s -> %s:%s, %d", flow.ForwardIP, flow.ForwardPort, flow.BackwardIP, flow.BackwardPort, flow.TransportProtocol)
}

// GetFlowMap return the flow map
func GetFlowMap() map[string]*Flow {
	return flowMap
}

// GetTotal return total count
func GetTotal() int {
	return total
}

// GetFlowPayloads Returns the flow payloads whose load exceeds the threshold
func GetFlowPayloads() []string {
	return flowPayloads
}

func ResetMap() {
	total += len(flowMap)
	flowMap = map[string]*Flow{}
}

// Clean reset flowMap flowPayloads flowContentPayloads
func Clean() {
	flowMap = map[string]*Flow{}
	flowPayloads = []string{}
}

// SetFlow set the flow map
func SetFlow(setMap map[string]*Flow, setFlowPayloads []string) {
	flowMap = setMap
	flowPayloads = setFlowPayloads
}
