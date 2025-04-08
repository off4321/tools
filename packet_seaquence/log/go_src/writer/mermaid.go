// Package writer provides functionality to write sequence diagrams
package writer

import (
	"fmt"
	"strings"

	"app/go_src/models"
)

// WriteMermaid provides functionality to generate Mermaid sequence diagrams
type WriteMermaid struct {
	Data *models.PacketSequenceData
}

// NewWriteMermaid creates a new WriteMermaid instance
func NewWriteMermaid(data *models.PacketSequenceData) *WriteMermaid {
	return &WriteMermaid{
		Data: data,
	}
}

// Generate generates a Mermaid sequence diagram
func (w *WriteMermaid) Generate() string {
	// Start building the Mermaid diagram
	var sb strings.Builder
	
	// Add the Mermaid header
	sb.WriteString("sequenceDiagram\n")
	
	// Get all nodes (participants)
	nodes := w.Data.GetNodes()
	for _, node := range nodes {
		// Sanitize node name for Mermaid syntax (replace dots with underscores)
		sanitizedNode := strings.ReplaceAll(node, ".", "_")
		sb.WriteString(fmt.Sprintf("    participant %s\n", sanitizedNode))
	}
	
	// Add all packets as arrows
	packets := w.Data.GetPackets()
	for _, packet := range packets {
		// Sanitize source and destination for Mermaid syntax
		source := strings.ReplaceAll(packet.Source, ".", "_")
		dest := strings.ReplaceAll(packet.Dest, ".", "_")
		
		// Build the message text
		message := w.BuildMessage(packet)
		
		// Add the arrow
		sb.WriteString(fmt.Sprintf("    %s->>%s: %s\n", source, dest, message))
	}
	
	return sb.String()
}

// BuildMessage builds a message for a packet
func (w *WriteMermaid) BuildMessage(packet models.PacketInfo) string {
	// デバッグ: パケット情報の確認
	fmt.Printf("DEBUG: BuildMessage - Protocol: %s, ProtocolNameFromPacket: %s\n", 
		packet.Protocol, packet.ProtocolNameFromPacket)
	fmt.Printf("DEBUG: BuildMessage - Info keys: %v\n", packet.Info)
	
	// Start with the highest layer protocol
	var message string

	// Check for higher-layer protocols first
	if httpInfo, ok := extractProtocolInfo(packet.Info, "http_info"); ok {
		method, _ := httpInfo["method"].(string)
		uri, _ := httpInfo["uri"].(string)
		if method != "" && uri != "" {
			message = fmt.Sprintf("HTTP %s %s", method, uri)
			return message
		}
	} else if dnsInfo, ok := extractProtocolInfo(packet.Info, "dns_info"); ok {
		queryName, _ := dnsInfo["query_name"].(string)
		if queryName != "" {
			message = fmt.Sprintf("DNS %s", queryName)
			return message
		}
	} else if tcpInfo, ok := extractProtocolInfo(packet.Info, "tcp_info"); ok {
		srcPort, _ := tcpInfo["src_port"].(string)
		dstPort, _ := tcpInfo["dst_port"].(string)
		if srcPort != "" && dstPort != "" {
			message = fmt.Sprintf("TCP %s→%s", srcPort, dstPort)
			return message
		}
	} else if udpInfo, ok := extractProtocolInfo(packet.Info, "udp_info"); ok {
		srcPort, _ := udpInfo["src_port"].(string)
		dstPort, _ := udpInfo["dst_port"].(string)
		if srcPort != "" && dstPort != "" {
			// NTPチェック - UDPポート123を使用するケース
			if (srcPort == "123" || dstPort == "123") {
				message = fmt.Sprintf("NTP %s→%s", srcPort, dstPort)
				return message
			}
			message = fmt.Sprintf("UDP %s→%s", srcPort, dstPort)
			return message
		}
	} else if sctpInfo, ok := extractProtocolInfo(packet.Info, "sctp_info"); ok {
		srcPort, _ := sctpInfo["src_port"].(float64)
		dstPort, _ := sctpInfo["dst_port"].(float64)
		if srcPort != 0 && dstPort != 0 {
			message = fmt.Sprintf("SCTP %v→%v", srcPort, dstPort)
			return message
		}
	} else if icmpInfo, ok := extractProtocolInfo(packet.Info, "icmp_info"); ok {
		typeName, _ := icmpInfo["type_name"].(string)
		if typeName != "" {
			message = fmt.Sprintf("ICMP %s", typeName)
			return message
		}
	} else if arpInfo, ok := extractProtocolInfo(packet.Info, "arp_info"); ok {
		// ARPプロトコルの場合はIPアドレス情報を表示する
		srcIP, srcOk := arpInfo["src_ip"].(string)
		dstIP, dstOk := arpInfo["dst_ip"].(string)
		
		// IPアドレス情報があれば表示
		if srcOk && dstOk && srcIP != "" && dstIP != "" {
			message = fmt.Sprintf("ARP %s→%s", srcIP, dstIP)
			return message
		}
		
		// 片方のIPだけある場合
		if srcOk && srcIP != "" {
			message = fmt.Sprintf("ARP from %s", srcIP)
			return message
		}
		if dstOk && dstIP != "" {
			message = fmt.Sprintf("ARP to %s", dstIP)
			return message
		}
		
		// IPがどちらもなければ操作を表示
		operation, opOk := arpInfo["operation"].(string)
		if opOk && operation != "" {
			message = fmt.Sprintf("ARP %s", operation)
			return message
		}
	}
	
	// IPv4情報があればIPアドレスを表示
	if ipv4Info, ok := extractProtocolInfo(packet.Info, "ipv4_info"); ok {
		srcIP, srcOk := ipv4Info["src"].(string)
		dstIP, dstOk := ipv4Info["dst"].(string)
		
		if srcOk && dstOk && srcIP != "" && dstIP != "" {
			// 既存のメッセージがあれば、IPアドレス情報を追加
			if message != "" {
				message = fmt.Sprintf("%s (%s→%s)", message, srcIP, dstIP)
			} else {
				message = fmt.Sprintf("IPv4 %s→%s", srcIP, dstIP)
			}
			return message
		}
	}
	
	// Default to the packet's protocol if no higher-layer protocol is found
	if message == "" {
		message = packet.Protocol
		
		// ソースとデスティネーションの情報を追加
		if packet.Source != "Unknown" && packet.Dest != "Remote" {
			message = fmt.Sprintf("%s (%s→%s)", message, packet.Source, packet.Dest)
		}
	}
	
	return message
}

// Helper function to extract protocol info from packet info
func extractProtocolInfo(info map[string]interface{}, key string) (map[string]interface{}, bool) {
	if infoValue, ok := info[key]; ok {
		if infoMap, ok := infoValue.(map[string]interface{}); ok {
			return infoMap, true
		}
	}
	return nil, false
}