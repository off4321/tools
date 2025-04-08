// Package extractor provides functionality to extract data from packets
package extractor

import (
	"fmt"
	"time"
	
	"app/go_src/discriminator"
	"app/go_src/models"
	"app/go_src/protocol_analyzer"
)

// ExtractData provides functionality to extract data from packets
type ExtractData struct {
	Packets    []protocol_analyzer.Packet
	MaxEntries int
}

// NewExtractData creates a new ExtractData instance
func NewExtractData(packets []protocol_analyzer.Packet, maxEntries int) *ExtractData {
	return &ExtractData{
		Packets:    packets,
		MaxEntries: maxEntries,
	}
}

// Extract extracts packet data and returns a PacketSequenceData
func (e *ExtractData) Extract() *models.PacketSequenceData {
	sequenceData := models.NewPacketSequenceData()
	
	// Process only up to MaxEntries packets
	maxPackets := len(e.Packets)
	if e.MaxEntries > 0 && e.MaxEntries < maxPackets {
		maxPackets = e.MaxEntries
	}
	
	// Extract data from each packet
	for i := 0; i < maxPackets; i++ {
		packet := e.Packets[i]
		packetInfo, err := e.extractPacketData(packet)
		if err != nil {
			// Optionally log error
			fmt.Printf("Error extracting packet data: %v\n", err)
			continue
		}
		
		sequenceData.AddPacket(packetInfo)
	}
	
	return sequenceData
}

// extractPacketData extracts data from a single packet
func (e *ExtractData) extractPacketData(packet protocol_analyzer.Packet) (models.PacketInfo, error) {
	// Create a new discriminator for protocol detection
	protocolDiscriminator := discriminator.NewDiscriminateProtocol(packet)
	protocolInfo := protocolDiscriminator.Discriminate()
	
	// protocolInfoがnilの場合でも基本情報を作成する
	if protocolInfo == nil {
		// 基本的な情報のマップを作成
		protocolInfo = make(map[string]interface{})
		protocolInfo["protocol_name"] = "Unknown"
		
		// デバッグ情報
		fmt.Println("DEBUG: Creating basic info for undiscriminated protocol")
	}
	
	// Extract timestamp
	timestamp := e.extractTimestamp(packet)
	
	// Create empty IP info
	ipInfo := make(map[string]interface{})
	icmpInfo := make(map[string]interface{})
	arpInfo := make(map[string]interface{})
	
	// 未サポートプロトコル情報
	unsupportedInfo := make(map[string]interface{})
	
	// Get IP info if available
	if ipv4InfoRaw, ok := protocolInfo["ipv4_info"]; ok {
		if ipv4Info, ok := ipv4InfoRaw.(map[string]interface{}); ok {
			ipInfo = ipv4Info
		}
	}
	
	// Get ICMP info if available
	if icmpInfoRaw, ok := protocolInfo["icmp_info"]; ok {
		if icmpInfo_, ok := icmpInfoRaw.(map[string]interface{}); ok {
			icmpInfo = icmpInfo_
		}
	}
	
	// Get ARP info if available
	if arpInfoRaw, ok := protocolInfo["arp_info"]; ok {
		if arpInfo_, ok := arpInfoRaw.(map[string]interface{}); ok {
			arpInfo = arpInfo_
		}
	}
	
	// Get unsupported protocol info if available
	if unsupportedInfoRaw, ok := protocolInfo["unsupported_info"]; ok {
		if unsupportedInfo_, ok := unsupportedInfoRaw.(map[string]interface{}); ok {
			unsupportedInfo = unsupportedInfo_
		}
	}
	
	// Determine protocol and get source and destination
	protocol := e.determineProtocol(packet, protocolInfo, ipInfo, icmpInfo)
	src, dst := e.getSrcDst(ipInfo, protocolInfo, arpInfo, unsupportedInfo)
	
	// ソースとデスティネーションが取得できない場合は
	// 代替値を設定
	if src == "?" {
		src = "Unknown"
	}
	if dst == "?" {
		dst = "Remote"
	}
	
	// Create the packet info with ProtocolNameFromPacket set to the detected protocol
	packetInfo := models.PacketInfo{
		Time:                 timestamp,
		Source:               src,
		Dest:                 dst,
		Protocol:             protocol,
		ProtocolNameFromPacket: protocol, // パケットから直接取得したプロトコル名を保存
		Info:                 protocolInfo,
	}
	
	return packetInfo, nil
}

// extractTimestamp extracts the timestamp from a packet
func (e *ExtractData) extractTimestamp(packet protocol_analyzer.Packet) time.Time {
	// Try to cast the packet to SimplePacket to get timestamp
	if simplePacket, ok := packet.(interface{ GetTimestamp() time.Time }); ok {
		return simplePacket.GetTimestamp()
	}
	
	// Default: return current time
	return time.Now()
}

// determineProtocol determines the protocol name from protocol info
func (e *ExtractData) determineProtocol(packet protocol_analyzer.Packet, protocolInfo map[string]interface{}, ipInfo map[string]interface{}, icmpInfo map[string]interface{}) string {
	// 最も高いレイヤーのプロトコルから順にチェック
	
	// HTTPチェック
	if _, ok := protocolInfo["http_info"]; ok {
		return "HTTP"
	}
	
	// DNSチェック
	if _, ok := protocolInfo["dns_info"]; ok {
		return "DNS"
	}
	
	// NTPチェック - UDPポート123を使用するケース
	if udpInfo, ok := protocolInfo["udp_info"].(map[string]interface{}); ok {
		srcPort, _ := udpInfo["src_port"].(string)
		dstPort, _ := udpInfo["dst_port"].(string)
		if srcPort == "123" || dstPort == "123" {
			return "NTP"
		}
	}
	
	// SCTPチェック
	if _, ok := protocolInfo["sctp_info"]; ok {
		return "SCTP"
	}
	
	// TCPチェック
	if _, ok := protocolInfo["tcp_info"]; ok {
		return "TCP"
	}
	
	// UDPチェック
	if _, ok := protocolInfo["udp_info"]; ok {
		return "UDP"
	}
	
	// ICMPチェック
	if _, ok := protocolInfo["icmp_info"]; ok {
		// Get ICMP type for more specific protocol name
		if len(icmpInfo) > 0 {
			if typeName, ok := icmpInfo["type_name"].(string); ok {
				if typeName == "Echo Request" {
					return "ICMP ECHO_REQUEST"
				} else if typeName == "Echo Reply" {
					return "ICMP ECHO_REPLY"
				} else {
					return "ICMP " + typeName
				}
			}
		}
		return "ICMP"
	}
	
	// ARPチェック
	if _, ok := protocolInfo["arp_info"]; ok {
		// Check ARP operation
		if arpInfo, ok := protocolInfo["arp_info"].(map[string]interface{}); ok {
			if opcode, ok := arpInfo["opcode"].(string); ok {
				if opcode == "1" {
					return "ARP REQUEST"
				} else if opcode == "2" {
					return "ARP REPLY"
				}
			}
		}
		return "ARP"
	}
	
	// X.25プロトコルチェック
	if _, ok := protocolInfo["x25_info"]; ok {
		return "X.25"
	}
	
	// IPARSプロトコルチェック - レイヤーに直接表示されている場合
	if packet.LayerExists("ipars") {
		return "IPARS"
	}
	
	// LAPBプロトコルチェック
	if packet.LayerExists("lapb") {
		return "LAPB"
	}
	
	// 未サポートのプロトコル情報がある場合
	if unsupportedInfo, ok := protocolInfo["unsupported_info"].(map[string]interface{}); ok {
		if protocolName, ok := unsupportedInfo["protocol_name"].(string); ok {
			return protocolName
		}
	}
	
	// プロトコル名が直接指定されている場合
	if protocolName, ok := protocolInfo["protocol_name"].(string); ok {
		if protocolName != "Unknown" {
			return protocolName
		}
	}
	
	// IPv4チェック (最も低いレイヤー)
	if _, ok := protocolInfo["ipv4_info"]; ok {
		return "IPv4"
	}
	
	// generic_infoからプロトコル名を取得
	if genericInfo, ok := protocolInfo["generic_info"].(map[string]interface{}); ok {
		if name, ok := genericInfo["name"].(string); ok && name != "Unknown" {
			return name
		}
	}
	
	// パケットから直接プロトコル名を取得（最後の手段）
	protocolName := packet.GetProtocolName()
	if protocolName != "" && protocolName != "Unknown" {
		return protocolName
	}
	
	// Default: unknown
	return "UNKNOWN"
}

// getSrcDst extracts source and destination from protocol info
func (e *ExtractData) getSrcDst(ipInfo map[string]interface{}, protocolInfo map[string]interface{}, arpInfo map[string]interface{}, unsupportedInfo map[string]interface{}) (string, string) {
	// Default values
	src := "?"
	dst := "?"
	
	// Try to get from IPv4 info
	if ipInfo != nil {
		// Check both 'src_ip' and 'src' fields (for backwards compatibility)
		if srcIP, ok := ipInfo["src_ip"].(string); ok {
			src = srcIP
		} else if srcIP, ok := ipInfo["src"].(string); ok {
			src = srcIP
		}
		
		// Check both 'dst_ip' and 'dst' fields (for backwards compatibility)
		if dstIP, ok := ipInfo["dst_ip"].(string); ok {
			dst = dstIP
		} else if dstIP, ok := ipInfo["dst"].(string); ok {
			dst = dstIP
		}
	}
	
	// If not found in IP, try ARP
	if src == "?" && dst == "?" && len(arpInfo) > 0 {
		if srcIP, ok := arpInfo["src_proto_ipv4"].(string); ok {
			src = srcIP
		}
		if dstIP, ok := arpInfo["dst_proto_ipv4"].(string); ok {
			dst = dstIP
		}
	}
	
	// If not found in IP or ARP, try to get from unsupported protocol info
	if src == "?" && dst == "?" && len(unsupportedInfo) > 0 {
		if srcAddr, ok := unsupportedInfo["src"].(string); ok {
			src = srcAddr
		}
		if dstAddr, ok := unsupportedInfo["dst"].(string); ok {
			dst = dstAddr
		}
	}
	
	// If still not found, try to get from protocol name
	if src == "?" && dst == "?" {
		if protocolName, ok := protocolInfo["protocol_name"].(string); ok {
			// LAPBやIPARS等の特定プロトコルの場合、デフォルト値を設定
			if protocolName == "LAPB" || protocolName == "IPARS" || 
			   protocolName == "lapb" || protocolName == "ipars" {
				src = "Local"
				dst = "Remote"
			}
		}
	}
	
	return src, dst
}