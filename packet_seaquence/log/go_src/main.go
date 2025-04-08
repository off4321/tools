// Package main provides a sample implementation of the packet analyzer
package main

import (
	"fmt"
	
	"app/go_src/discriminator"
	"app/go_src/extractor"
	"app/go_src/packet"
	"app/go_src/protocol_analyzer"
)

func main() {
	fmt.Println("Go言語パケット解析ツール")
	fmt.Println("========================")
	
	// サンプルパケットを作成
	packets := createSamplePackets()
	
	// 各パケットを解析
	for i, pkt := range packets {
		fmt.Printf("\nパケット #%d:\n", i+1)
		fmt.Println("-------------")
		
		// プロトコル判別
		protocolDiscriminator := discriminator.NewDiscriminateProtocol(pkt)
		protocolInfo := protocolDiscriminator.Discriminate()
		
		if protocolInfo == nil {
			fmt.Println("プロトコル情報を取得できませんでした")
			continue
		}
		
		// プロトコルアナライザを使用して各プロトコルの詳細情報を表示
		displayProtocolInfo(pkt, protocolInfo)
	}
	
	// 全パケットからデータを抽出
	fmt.Println("\n\nパケットシーケンスデータ:")
	fmt.Println("=====================")
	extractor := extractor.NewExtractData(packets, 10)
	sequenceData := extractor.Extract()
	
	// 結果を表示
	fmt.Printf("取得したパケット数: %d\n", len(sequenceData.GetPackets()))
	fmt.Printf("ノード数: %d\n", len(sequenceData.GetNodes()))
	fmt.Println("ノード一覧:", sequenceData.GetNodes())
	
	// 各パケットの基本情報を表示
	fmt.Println("\nパケット概要:")
	for i, pktInfo := range sequenceData.GetPackets() {
		fmt.Printf("#%d: %s -> %s (%s) at %v\n", 
			i+1, pktInfo.Source, pktInfo.Dest, pktInfo.Protocol, pktInfo.Time)
	}
}

// createSamplePackets creates sample packets for demonstration
func createSamplePackets() []protocol_analyzer.Packet {
	packets := make([]protocol_analyzer.Packet, 0)
	
	// サンプル1: IPとTCPを含むパケット (HTTP)
	packet1 := packet.NewSimplePacket()
	
	// イーサネットレイヤー
	packet1.AddLayer("eth", map[string]interface{}{
		"src": "00:11:22:33:44:55",
		"dst": "aa:bb:cc:dd:ee:ff",
		"type": "IPv4",
	})
	
	// IPレイヤー
	packet1.AddLayer("ip", map[string]interface{}{
		"src": "192.168.1.100",
		"dst": "93.184.216.34",
		"ttl": "64",
		"flags": "0x02",
		"len": "52",
		"protocol": "6", // TCP
	})
	
	// TCPレイヤー
	packet1.AddLayer("tcp", map[string]interface{}{
		"srcport": "49152",
		"dstport": "80",
		"seq": "1000",
		"seq_raw": "1000",
		"ack": "2000",
		"ack_raw": "2000",
		"flags": "0x18", // PSH, ACK
		"window_size": "64240",
	})
	
	// HTTPレイヤー
	packet1.AddLayer("http", map[string]interface{}{
		"request_method": "GET",
		"request_uri": "/index.html",
		"host": "example.com",
		"request_line": "GET /index.html HTTP/1.1",
	})
	
	packets = append(packets, packet1)
	
	// サンプル2: IPとUDPを含むパケット (DNS)
	packet2 := packet.NewSimplePacket()
	
	// イーサネットレイヤー
	packet2.AddLayer("eth", map[string]interface{}{
		"src": "00:11:22:33:44:55",
		"dst": "aa:bb:cc:dd:ee:ff",
		"type": "IPv4",
	})
	
	// IPレイヤー
	packet2.AddLayer("ip", map[string]interface{}{
		"src": "192.168.1.100",
		"dst": "8.8.8.8",
		"ttl": "64",
		"flags": "0x00",
		"len": "60",
		"protocol": "17", // UDP
	})
	
	// UDPレイヤー
	packet2.AddLayer("udp", map[string]interface{}{
		"srcport": "53568",
		"dstport": "53",
		"length": "32",
	})
	
	// DNSレイヤー
	packet2.AddLayer("dns", map[string]interface{}{
		"qry_name": "example.com",
		"qry_type": "A",
	})
	
	packets = append(packets, packet2)
	
	// サンプル3: ARPパケット
	packet3 := packet.NewSimplePacket()
	
	// イーサネットレイヤー
	packet3.AddLayer("eth", map[string]interface{}{
		"src": "00:11:22:33:44:55",
		"dst": "ff:ff:ff:ff:ff:ff", // ブロードキャスト
		"type": "ARP",
	})
	
	// ARPレイヤー
	packet3.AddLayer("arp", map[string]interface{}{
		"opcode": "1", // リクエスト
		"src_hw_mac": "00:11:22:33:44:55",
		"src_proto_ipv4": "192.168.1.100",
		"dst_hw_mac": "00:00:00:00:00:00",
		"dst_proto_ipv4": "192.168.1.1",
	})
	
	packets = append(packets, packet3)
	
	// サンプル4: IPとTCPを含むパケット (HTTPS)
	packet4 := packet.NewSimplePacket()
	
	// イーサネットレイヤー
	packet4.AddLayer("eth", map[string]interface{}{
		"src": "00:11:22:33:44:55",
		"dst": "aa:bb:cc:dd:ee:ff",
		"type": "IPv4",
	})
	
	// IPレイヤー
	packet4.AddLayer("ip", map[string]interface{}{
		"src": "192.168.1.100",
		"dst": "93.184.216.34",
		"ttl": "64",
		"flags": "0x02",
		"len": "52",
		"protocol": "6", // TCP
	})
	
	// TCPレイヤー
	packet4.AddLayer("tcp", map[string]interface{}{
		"srcport": "49152",
		"dstport": "443",
		"seq": "1000",
		"seq_raw": "1000",
		"ack": "2000",
		"ack_raw": "2000",
		"flags": "0x18", // PSH, ACK
		"window_size": "64240",
	})
	
	// TLSレイヤー
	packet4.AddLayer("tls", map[string]interface{}{
		"record_type": "Application Data",
		"version": "TLS 1.2",
	})
	
	packets = append(packets, packet4)
	
	return packets
}

// displayProtocolInfo displays detailed information about each protocol in the packet
func displayProtocolInfo(pkt protocol_analyzer.Packet, protocolInfo map[string]interface{}) {
	// 利用可能なプロトコルの一覧を表示
	var protocols []string
	for key := range protocolInfo {
		if key != "protocol_name" {
			protocols = append(protocols, key)
		}
	}
	
	fmt.Println("検出されたプロトコル:")
	
	// IPv4: 他の上位プロトコル解析結果が無い場合のみ表示
	if !(protocolInfo["tcp_info"] != nil ||
		protocolInfo["udp_info"] != nil ||
		protocolInfo["sctp_info"] != nil ||
		protocolInfo["dns_info"] != nil ||
		protocolInfo["http_info"] != nil ||
		protocolInfo["https_info"] != nil ||
		protocolInfo["x25_info"] != nil) {
		analyzer := protocol_analyzer.NewAnalyzeIPv4(pkt)
		fmt.Printf("- %s\n", analyzer.GetDisplayInfo())
	}
	
	// ARP
	if _, ok := protocolInfo["arp_info"]; ok {
		analyzer := protocol_analyzer.NewAnalyzeARP(pkt)
		fmt.Printf("- %s\n", analyzer.GetDisplayInfo())
	}
	
	// TCP
	if _, ok := protocolInfo["tcp_info"]; ok {
		analyzer := protocol_analyzer.NewAnalyzeTCP(pkt)
		fmt.Printf("- %s\n", analyzer.GetDisplayInfo())
	}
	
	// UDP
	if _, ok := protocolInfo["udp_info"]; ok {
		analyzer := protocol_analyzer.NewAnalyzeUDP(pkt)
		fmt.Printf("- %s\n", analyzer.GetDisplayInfo())
	}
	
	// SCTP
	if _, ok := protocolInfo["sctp_info"]; ok {
		analyzer := protocol_analyzer.NewAnalyzeSCTP(pkt)
		fmt.Printf("- %s\n", analyzer.GetDisplayInfo())
	}
	
	// DNS
	if _, ok := protocolInfo["dns_info"]; ok {
		analyzer := protocol_analyzer.NewAnalyzeDNS(pkt)
		fmt.Printf("- %s\n", analyzer.GetDisplayInfo())
	}
	
	// HTTP
	if _, ok := protocolInfo["http_info"]; ok {
		analyzer := protocol_analyzer.NewAnalyzeHTTP(pkt)
		fmt.Printf("- %s\n", analyzer.GetDisplayInfo())
	}
	
	// HTTPS
	if _, ok := protocolInfo["https_info"]; ok {
		analyzer := protocol_analyzer.NewAnalyzeHTTPS(pkt)
		fmt.Printf("- %s\n", analyzer.GetDisplayInfo())
	}
	
	// X.25
	if _, ok := protocolInfo["x25_info"]; ok {
		analyzer := protocol_analyzer.NewAnalyzeX25(pkt)
		fmt.Printf("- %s\n", analyzer.GetDisplayInfo())
	}
	
	// 認識されなかったプロトコル
	if protocolName, ok := protocolInfo["protocol_name"].(string); ok {
		if len(protocolInfo) == 1 {
			analyzer := protocol_analyzer.NewAnalyzeUnsupportedProtocol(pkt, protocolName)
			fmt.Printf("- %s\n", analyzer.GetDisplayInfo())
		}
	}
}