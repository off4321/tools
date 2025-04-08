// Package models defines data structures for packet analysis
package models

import (
	"time"
)

// PacketInfo represents information about a single packet
type PacketInfo struct {
	Time                 time.Time              `json:"time"`
	Source               string                 `json:"src"`
	Dest                 string                 `json:"dst"`
	Protocol             string                 `json:"protocol"`
	ProtocolNameFromPacket string               `json:"protocol_name_from_packet"`
	Info                 map[string]interface{} `json:"info"`
}

// PacketSequenceData stores information about a sequence of packets
type PacketSequenceData struct {
	Packets []PacketInfo `json:"packets"`
	Nodes   map[string]struct{} `json:"nodes"`
}

// NewPacketSequenceData creates a new empty PacketSequenceData
func NewPacketSequenceData() *PacketSequenceData {
	return &PacketSequenceData{
		Packets: make([]PacketInfo, 0),
		Nodes:   make(map[string]struct{}),
	}
}

// AddPacket adds a packet to the sequence data
func (p *PacketSequenceData) AddPacket(packetInfo PacketInfo) {
	// 1. まず ProtocolNameFromPacket の値を確認
	if packetInfo.ProtocolNameFromPacket != "" && packetInfo.ProtocolNameFromPacket != "UNKNOWN" {
		// 検出されたプロトコル名がある場合はそれを使用
		packetInfo.Protocol = packetInfo.ProtocolNameFromPacket
	} else {
		// 2. ProtocolNameFromPacket が利用できない場合、従来のロジックで判定
		
		// 情報キーを確認して最高レイヤーのプロトコルを判断
		if _, ok := packetInfo.Info["http_info"]; ok {
			packetInfo.Protocol = "HTTP"
		} else if _, ok := packetInfo.Info["dns_info"]; ok {
			packetInfo.Protocol = "DNS"
		} else if _, ok := packetInfo.Info["sctp_info"]; ok {
			packetInfo.Protocol = "SCTP"
		} else if _, ok := packetInfo.Info["tcp_info"]; ok {
			packetInfo.Protocol = "TCP"
		} else if udpInfo, ok := packetInfo.Info["udp_info"].(map[string]interface{}); ok {
			// NTP判定 - UDPポート123を使用するケース
			srcPort, srcOk := udpInfo["src_port"].(string)
			dstPort, dstOk := udpInfo["dst_port"].(string)
			
			if (srcOk && srcPort == "123") || (dstOk && dstPort == "123") {
				packetInfo.Protocol = "NTP"
			} else {
				// その他のUDP
				packetInfo.Protocol = "UDP"
			}
		} else if _, ok := packetInfo.Info["icmp_info"]; ok {
			packetInfo.Protocol = "ICMP"
		} else if _, ok := packetInfo.Info["arp_info"]; ok {
			packetInfo.Protocol = "ARP"
		} else if _, ok := packetInfo.Info["ipv4_info"]; ok {
			packetInfo.Protocol = "IPv4"
		} else {
			// デフォルト設定
			packetInfo.Protocol = "UNKNOWN"
		}
	}

	// パケットをシーケンスに追加
	p.Packets = append(p.Packets, packetInfo)
	p.Nodes[packetInfo.Source] = struct{}{}
	p.Nodes[packetInfo.Dest] = struct{}{}
}

// GetPackets returns all packets in the sequence
func (p *PacketSequenceData) GetPackets() []PacketInfo {
	return p.Packets
}

// GetNodes returns all unique nodes (IP addresses, etc.) in the sequence
func (p *PacketSequenceData) GetNodes() []string {
	nodes := make([]string, 0, len(p.Nodes))
	for node := range p.Nodes {
		nodes = append(nodes, node)
	}
	return nodes
}

// CopyFirstN returns a new PacketSequenceData containing only the first n packets
func (p *PacketSequenceData) CopyFirstN(n int) *PacketSequenceData {
	newData := NewPacketSequenceData()
	
	for i, packet := range p.Packets {
		if i >= n {
			break
		}
		newData.AddPacket(packet)
	}
	
	return newData
}