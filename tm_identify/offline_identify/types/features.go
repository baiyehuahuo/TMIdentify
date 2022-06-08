package types

import (
	"github.com/google/gopacket/layers"
	"log"
	"math"
	"reflect"
)

type Features struct {
	// todo 可以考虑双向一起算的 还有偏斜 峰度等

	// The basic unit of time is Microseconds
	minFiat  int64   // Minimum of forward inter-arrival time
	meanFiat float64 // Mean of forward inter-arrival time
	maxFiat  int64   // Maximum of forward inter-arrival time
	stdFiat  float64 // Standard deviation of forward inter-arrival times

	minBiat  int64   // Minimum of backward inter-arrival time
	meanBiat float64 // Mean backward inter-arrival time
	maxBiat  int64   // Maximum of backward inter-arrival time
	stdBiat  float64 // Standard deviation of backward inter-arrival times

	minFpkt  int64   // Minimum of forward packet length
	meanFpkt float64 // Mean of forward packet length
	maxFpkt  int64   // Maximum of forward packet length
	stdFpkt  float64 // Standard deviation of forward packet length

	minBpkt  int64   // Minimum of backward packet length
	meanBpkt float64 // Mean of backward packet length
	maxBpkt  int64   // Maximum of backward packet length
	stdBpkt  float64 // Standard deviation of backward packet length

	protocol layers.IPProtocol
	duration int64 // Total duration
	hasRTCP  int   // Whether the packets contain the RTCP protocol

	fPackets int // Number of packets in forward direction
	fBytes   int // Number of bytes in forward direction

	bPackets int // Number of packets in backward direction
	bBytes   int // Number of bytes in backward direction
}

func getIATs(feature []PacketFeature, index []int) (iats []int64) {
	length := len(index)
	if length < 2 {
		return
	}
	iats = make([]int64, 0, length-1)
	for i := 1; i < length; i++ {
		iats = append(iats, feature[index[i]].Timestamp.Sub(feature[index[i-1]].Timestamp).Microseconds())
	}
	return
}

func getPayloadsLength(feature []PacketFeature, index []int) (payloadsLength []int64) {
	length := len(index)
	payloadsLength = make([]int64, 0, length)
	for i := range index {
		payloadsLength = append(payloadsLength, int64(feature[index[i]].PayloadLength))
	}
	return
}

func getInt64MinMeanMaxStd(slice []int64) (min int64, mean float64, max int64, std float64) {
	var total int64
	min, max = slice[0], slice[0]
	length := int64(len(slice))
	for i := range slice {
		if slice[i] < min {
			min = slice[i]
		}
		if slice[i] > max {
			max = slice[i]
		}
		total += slice[i]
	}
	mean = float64(total) / float64(length)
	var variance, t float64
	for i := range slice {
		t = float64(slice[i]) - mean
		variance += t * t
	}
	std = math.Sqrt(variance / float64(length))
	return
}

func getPacketAndByteNum(feature []PacketFeature, index []int) (packets int, bytes int) {
	packets = len(index)
	for i := range index {
		bytes += feature[index[i]].PayloadLength
	}
	return
}

// CalculateFeaturesFromFlow Compute identifying features for a flow to identify
func CalculateFeaturesFromFlow(flow *Flow) Features {
	var features Features

	if iats := getIATs(flow.PacketsFeature, flow.ForwardPackets); len(iats) > 0 {
		features.minFiat, features.meanFiat, features.maxFiat, features.stdFiat = getInt64MinMeanMaxStd(iats)
	}

	if iats := getIATs(flow.PacketsFeature, flow.BackwardPackets); len(iats) > 0 {
		features.minBiat, features.meanBiat, features.maxBiat, features.stdBiat = getInt64MinMeanMaxStd(iats)
	}

	if payloadsLength := getPayloadsLength(flow.PacketsFeature, flow.ForwardPackets); len(payloadsLength) > 0 {
		features.minFpkt, features.meanFpkt, features.maxFpkt, features.stdFpkt = getInt64MinMeanMaxStd(payloadsLength)
	}

	if payloadsLength := getPayloadsLength(flow.PacketsFeature, flow.BackwardPackets); len(payloadsLength) > 0 {
		features.minBpkt, features.meanBpkt, features.maxBpkt, features.stdBpkt = getInt64MinMeanMaxStd(payloadsLength)
	}

	features.protocol = flow.TransportProtocol
	features.duration = flow.PacketsFeature[len(flow.PacketsFeature)-1].Timestamp.Sub(flow.PacketsFeature[0].Timestamp).Microseconds()
	features.hasRTCP = flow.HasRTCP

	features.fPackets, features.fBytes = getPacketAndByteNum(flow.PacketsFeature, flow.ForwardPackets)
	features.bPackets, features.bBytes = getPacketAndByteNum(flow.PacketsFeature, flow.BackwardPackets)
	return features
}

func convertFeaturesToFloat64(features Features) []float64 {
	v := reflect.ValueOf(features)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	fieldNum := v.NumField()
	result := make([]float64, 0, fieldNum)
	var value reflect.Value
	for i := 0; i < fieldNum; i++ {
		value = v.Field(i)
		switch value.Kind().String() {
		case "uint8":
			result = append(result, float64(value.Uint()))
		case "int", "int64":
			result = append(result, float64(value.Int()))
		case "float64":
			result = append(result, value.Float())
		default:
			log.Fatalf("filed: %v", value.Kind().String())
		}
	}
	return result
}
