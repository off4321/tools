// Package protocol_analyzer provides ARP protocol analyzer
package protocol_analyzer

import (
	"fmt"
	"strconv"
	"strings"
)

// AnalyzeARP implements the ProtocolAnalyzer interface for ARP
type AnalyzeARP struct {
	Pkt Packet
}

// NewAnalyzeARP creates a new ARP analyzer
func NewAnalyzeARP(packet Packet) *AnalyzeARP {
	return &AnalyzeARP{
		Pkt: packet,
	}
}

// Analyze extracts ARP information from a packet
func (a *AnalyzeARP) Analyze() map[string]interface{} {
	if !a.Pkt.LayerExists("arp") {
		return nil
	}

	arpInfo := make(map[string]interface{})
	
	// ARP operation code
	if opcodeVal, ok := a.Pkt.GetField("arp", "opcode"); ok {
		opcode, err := strconv.Atoi(fmt.Sprintf("%v", opcodeVal))
		if err == nil {
			arpInfo["operation_code"] = opcode
			arpInfo["operation"] = a.getARPOperation(opcode)
		}
	} else if operationVal, ok := a.Pkt.GetField("arp", "operation"); ok {
		operation, err := strconv.Atoi(fmt.Sprintf("%v", operationVal))
		if err == nil {
			arpInfo["operation_code"] = operation
			arpInfo["operation"] = a.getARPOperation(operation)
		}
	}
	
	// Source MAC address
	if srcMac, ok := a.Pkt.GetField("arp", "src_hw_mac"); ok {
		arpInfo["src_mac"] = srcMac
	} else if srcMac, ok := a.Pkt.GetField("arp", "src.hw_mac"); ok {
		arpInfo["src_mac"] = srcMac
	}
	
	// Source IP address
	if srcIP, ok := a.Pkt.GetField("arp", "src_proto_ipv4"); ok {
		arpInfo["src_ip"] = srcIP
	} else if srcIP, ok := a.Pkt.GetField("arp", "src.proto_ipv4"); ok {
		arpInfo["src_ip"] = srcIP
	}
	
	// Destination MAC address
	if dstMac, ok := a.Pkt.GetField("arp", "dst_hw_mac"); ok {
		arpInfo["dst_mac"] = dstMac
	} else if dstMac, ok := a.Pkt.GetField("arp", "dst.hw_mac"); ok {
		arpInfo["dst_mac"] = dstMac
	}
	
	// Destination IP address
	if dstIP, ok := a.Pkt.GetField("arp", "dst_proto_ipv4"); ok {
		arpInfo["dst_ip"] = dstIP
	} else if dstIP, ok := a.Pkt.GetField("arp", "dst.proto_ipv4"); ok {
		arpInfo["dst_ip"] = dstIP
	}
	
	return arpInfo
}

// GetDisplayInfo returns formatted ARP information
func (a *AnalyzeARP) GetDisplayInfo() string {
	info := a.Analyze()
	if info == nil {
		return "ARP情報なし"
	}
	
	operation, _ := info["operation"].(string)
	if operation == "" {
		operation = "不明"
	}
	
	srcMac, _ := info["src_mac"].(string)
	if srcMac == "" {
		srcMac = "不明"
	}
	
	srcIP, _ := info["src_ip"].(string)
	if srcIP == "" {
		srcIP = "不明"
	}
	
	dstMac, _ := info["dst_mac"].(string)
	if dstMac == "" {
		dstMac = "不明"
	}
	
	dstIP, _ := info["dst_ip"].(string)
	if dstIP == "" {
		dstIP = "不明"
	}
	
	if operation == "REQUEST" {
		return fmt.Sprintf("ARP Request: %s (%s) は %s のMACアドレスを要求", srcIP, srcMac, dstIP)
	} else if operation == "REPLY" {
		return fmt.Sprintf("ARP Reply: %s は %s です", srcIP, srcMac)
	} else {
		return fmt.Sprintf("ARP %s: %s (%s) -> %s (%s)", operation, srcIP, srcMac, dstIP, dstMac)
	}
}

// GetSummary returns a summary of ARP information
func (a *AnalyzeARP) GetSummary() string {
	info := a.Analyze()
	if info == nil {
		return "ARP"
	}
	
	var operation string
	if op, ok := info["operation"].(string); ok {
		// Split at hyphen and get first part (e.g., "REPLY-REQUEST" -> "REPLY")
		parts := strings.Split(op, "-")
		operation = parts[0]
	} else {
		operation = ""
	}
	
	srcIP, _ := info["src_ip"].(string)
	if srcIP == "" {
		srcIP = "?"
	}
	
	dstIP, _ := info["dst_ip"].(string)
	if dstIP == "" {
		dstIP = "?"
	}
	
	return fmt.Sprintf("ARP %s %s -> %s", operation, srcIP, dstIP)
}

// getARPOperation returns a human-readable ARP operation string
func (a *AnalyzeARP) getARPOperation(code int) string {
	arpOperations := map[int]string{
		1: "REQUEST",
		2: "REPLY",
		3: "RARP-REQUEST",
		4: "RARP-REPLY",
		5: "DRARP-REQUEST",
		6: "DRARP-REPLY",
		7: "DRARP-ERROR",
		8: "InARP-REQUEST",
		9: "InARP-REPLY",
	}
	
	if op, ok := arpOperations[code]; ok {
		return op
	}
	return fmt.Sprintf("UNKNOWN (%d)", code)
}