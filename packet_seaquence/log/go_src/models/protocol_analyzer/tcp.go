// Package protocol_analyzer provides TCP protocol analyzer
package protocol_analyzer

import (
	"fmt"
	"strings"
)

// AnalyzeTCP implements the ProtocolAnalyzer interface for TCP
type AnalyzeTCP struct {
	BaseProtocolAnalyzer
}

// NewAnalyzeTCP creates a new TCP analyzer
func NewAnalyzeTCP(packet Packet) *AnalyzeTCP {
	return &AnalyzeTCP{
		BaseProtocolAnalyzer: NewBaseProtocolAnalyzer(packet),
	}
}

// Analyze extracts TCP information from a packet
func (a *AnalyzeTCP) Analyze() map[string]interface{} {
	if !a.Pkt.LayerExists("tcp") {
		return nil
	}

	tcpInfo := make(map[string]interface{})
	
	// Source port
	if srcPort, ok := a.Pkt.GetField("tcp", "srcport"); ok {
		tcpInfo["src_port"] = srcPort
	}
	
	// Destination port
	if dstPort, ok := a.Pkt.GetField("tcp", "dstport"); ok {
		tcpInfo["dst_port"] = dstPort
	}
	
	// Sequence number
	if seq, ok := a.Pkt.GetField("tcp", "seq"); ok {
		tcpInfo["seq"] = seq
		if seqRaw, ok := a.Pkt.GetField("tcp", "seq_raw"); ok {
			tcpInfo["seq_raw"] = seqRaw
		}
	}
	
	// Acknowledgment number
	if ack, ok := a.Pkt.GetField("tcp", "ack"); ok {
		tcpInfo["ack"] = ack
		if ackRaw, ok := a.Pkt.GetField("tcp", "ack_raw"); ok {
			tcpInfo["ack_raw"] = ackRaw
		}
	}
	
	// Flags
	if flags, ok := a.Pkt.GetField("tcp", "flags"); ok {
		tcpInfo["flags"] = flags
		tcpInfo["flags_desc"] = a.getTcpFlagsDesc(fmt.Sprintf("%v", flags))
	}
	
	// Window size
	if windowSize, ok := a.Pkt.GetField("tcp", "window_size"); ok {
		tcpInfo["window_size"] = windowSize
	} else if window, ok := a.Pkt.GetField("tcp", "window"); ok {
		tcpInfo["window_size"] = window
	}
	
	// Checksum
	if checksum, ok := a.Pkt.GetField("tcp", "checksum"); ok {
		tcpInfo["checksum"] = checksum
	}
	
	// Urgent pointer
	if urgentPointer, ok := a.Pkt.GetField("tcp", "urgent_pointer"); ok {
		tcpInfo["urgent_pointer"] = urgentPointer
	}
	
	// Options
	if options, ok := a.Pkt.GetField("tcp", "options"); ok {
		tcpInfo["options"] = options
	}
	
	// Payload
	if payload, ok := a.Pkt.GetField("tcp", "payload"); ok {
		tcpInfo["payload"] = payload
	} else if payloadRaw, ok := a.Pkt.GetField("tcp", "payload_raw"); ok {
		tcpInfo["payload_raw"] = payloadRaw
	}
	
	// Stream number
	if stream, ok := a.Pkt.GetField("tcp", "stream"); ok {
		tcpInfo["stream"] = stream
	}
	
	return tcpInfo
}

// GetDisplayInfo returns formatted TCP information
func (a *AnalyzeTCP) GetDisplayInfo() string {
	info := a.Analyze()
	if info == nil {
		return "TCP情報なし"
	}
	
	srcPort := getOrDefault(info["src_port"], "不明")
	dstPort := getOrDefault(info["dst_port"], "不明")
	flags := getOrDefault(info["flags_desc"], "不明")
	seq := getOrDefault(info["seq"], "N/A")
	ack := getOrDefault(info["ack"], "N/A")
	
	portInfo := fmt.Sprintf("%s -> %s", srcPort, dstPort)
	flagInfo := fmt.Sprintf("[%s]", flags)
	
	if strings.Contains(flags, "SYN") && !strings.Contains(flags, "ACK") {
		return fmt.Sprintf("TCP %s %s 接続開始", portInfo, flagInfo)
	} else if strings.Contains(flags, "SYN") && strings.Contains(flags, "ACK") {
		return fmt.Sprintf("TCP %s %s 接続確認", portInfo, flagInfo)
	} else if strings.Contains(flags, "FIN") {
		return fmt.Sprintf("TCP %s %s 接続終了", portInfo, flagInfo)
	} else if strings.Contains(flags, "RST") {
		return fmt.Sprintf("TCP %s %s 接続リセット", portInfo, flagInfo)
	} else if strings.Contains(flags, "ACK") {
		return fmt.Sprintf("TCP %s %s 確認応答 (SEQ=%s, ACK=%s)", portInfo, flagInfo, seq, ack)
	} else {
		return fmt.Sprintf("TCP %s %s", portInfo, flagInfo)
	}
}

// GetSummary returns a summary of TCP information
func (a *AnalyzeTCP) GetSummary() string {
	info := a.Analyze()
	if info == nil {
		return "TCP"
	}
	
	srcPort := getOrDefault(info["src_port"], "")
	dstPort := getOrDefault(info["dst_port"], "")
	flags := getOrDefault(info["flags_desc"], "")
	
	serviceInfo := a.getServiceNameFromPort(dstPort)
	
	if serviceInfo != "" {
		return fmt.Sprintf("TCP %s -> %s (%s) %s", srcPort, dstPort, serviceInfo, flags)
	} else {
		return fmt.Sprintf("TCP %s -> %s %s", srcPort, dstPort, flags)
	}
}

// getTcpFlagsDesc returns a description of TCP flags
func (a *AnalyzeTCP) getTcpFlagsDesc(flags string) string {
	var descriptions []string
	
	if strings.Contains(flags, "SYN") {
		descriptions = append(descriptions, "SYN")
	}
	if strings.Contains(flags, "ACK") {
		descriptions = append(descriptions, "ACK")
	}
	if strings.Contains(flags, "FIN") {
		descriptions = append(descriptions, "FIN")
	}
	if strings.Contains(flags, "RST") {
		descriptions = append(descriptions, "RST")
	}
	if strings.Contains(flags, "PSH") {
		descriptions = append(descriptions, "PSH")
	}
	if strings.Contains(flags, "URG") {
		descriptions = append(descriptions, "URG")
	}
	if strings.Contains(flags, "ECE") {
		descriptions = append(descriptions, "ECE")
	}
	if strings.Contains(flags, "CWR") {
		descriptions = append(descriptions, "CWR")
	}
	if strings.Contains(flags, "NS") {
		descriptions = append(descriptions, "NS")
	}
	
	if len(descriptions) > 0 {
		return strings.Join(descriptions, ", ")
	}
	return "None"
}

// getServiceNameFromPort returns a service name for a given port
func (a *AnalyzeTCP) getServiceNameFromPort(port string) string {
	wellKnownPorts := map[string]string{
		"20": "FTP-DATA",
		"21": "FTP",
		"22": "SSH",
		"23": "Telnet",
		"25": "SMTP",
		"53": "DNS",
		"67": "DHCP-Server",
		"68": "DHCP-Client",
		"80": "HTTP",
		"110": "POP3",
		"143": "IMAP",
		"443": "HTTPS",
		"465": "SMTPS",
		"993": "IMAPS",
		"995": "POP3S",
		"3389": "RDP",
	}
	
	if service, ok := wellKnownPorts[port]; ok {
		return service
	}
	return ""
}