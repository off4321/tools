// Package discriminator provides protocol discrimination functionality
package discriminator

import (
	"strings"
	"app/go_src/protocol_analyzer"
)

// DiscriminateProtocol provides functionality to identify protocols in packets
type DiscriminateProtocol struct {
	Pkt protocol_analyzer.Packet
}

// NewDiscriminateProtocol creates a new DiscriminateProtocol instance
func NewDiscriminateProtocol(packet protocol_analyzer.Packet) *DiscriminateProtocol {
	return &DiscriminateProtocol{
		Pkt: packet,
	}
}

// Discriminate identifies protocol information from a packet
func (d *DiscriminateProtocol) Discriminate() map[string]interface{} {
	// Create available discriminator to get layer information
	available := make(map[string]interface{})
	
	// Get protocol name from the packet (tshark based implementation will have this)
	protocolName := d.Pkt.GetProtocolName()
	if protocolName == "" {
		protocolName = "Unknown"
	}
	
	// For tshark implementation, we directly extract layers from the packet
	layers := []string{}
	
	// Check for specific layers using LayerExists
	// tsharkでは通常、frame.protocolsとして全レイヤーのリストが得られる
	commonLayers := []string{
		"ip", "ipv6", "icmp", "icmpv6", "tcp", "udp", "sctp", "arp", 
		"dns", "http", "ssl", "tls", "x25", "ipars", "lapb", 
		"uts", "q933", "igmp", "ntp", "smtp",
	}
	
	for _, layer := range commonLayers {
		if d.Pkt.LayerExists(layer) {
			layers = append(layers, layer)
		}
	}
	
	// tsharkの実装ではSrcIP、DstIPが直接アクセス可能なことがある
	if srcGetter, ok := d.Pkt.(interface { GetSrcIP() string }); ok {
		available["src_ip"] = srcGetter.GetSrcIP()
	}
	
	if dstGetter, ok := d.Pkt.(interface { GetDstIP() string }); ok {
		available["dst_ip"] = dstGetter.GetDstIP()
	}
	
	// タイムスタンプも直接アクセス可能なことがある
	if tsGetter, ok := d.Pkt.(interface { GetTimestamp() string }); ok {
		available["timestamp"] = tsGetter.GetTimestamp()
	}
	
	available["layer_names"] = layers
	
	// Determine highest layer
	highestLayer := protocolName
	// レイヤーがある場合は最後のレイヤーを最上位とする (ただしプロトコル名が既に設定されている場合を優先)
	if len(layers) > 0 && (highestLayer == "Unknown" || highestLayer == "") {
		highestLayer = layers[len(layers)-1]
	}
	
	available["highest_layer"] = highestLayer
	
	// 基本情報を含むprotocolInfoを初期化
	protocolInfo := make(map[string]interface{})
	protocolInfo["protocol_name"] = highestLayer
	
	// レイヤーが検出されなかった場合でも基本情報を返す
	if len(layers) == 0 {
		// 基本的な情報を追加
		genericInfo := make(map[string]interface{})
		genericInfo["name"] = highestLayer
		
		// IPアドレス情報があれば追加
		if srcIP, ok := available["src_ip"].(string); ok && srcIP != "" {
			genericInfo["src"] = srcIP
		}
		
		if dstIP, ok := available["dst_ip"].(string); ok && dstIP != "" {
			genericInfo["dst"] = dstIP
		}
		
		protocolInfo["generic_info"] = genericInfo
		
		// サポートされていないプロトコルとして処理
		analyzer := protocol_analyzer.NewAnalyzeUnsupportedProtocol(d.Pkt, highestLayer)
		protocolInfo["unsupported_info"] = analyzer.Analyze()
		
		return protocolInfo
	}
	
	// Extract protocol details using appropriate analyzers
	d.extractProtocolDetails(protocolInfo, available)
	
	return protocolInfo
}

// extractProtocolDetails analyzes each protocol layer using respective analyzers
func (d *DiscriminateProtocol) extractProtocolDetails(protocolInfo map[string]interface{}, available map[string]interface{}) map[string]interface{} {
	// Get layer names from available information
	layersInterface, ok := available["layer_names"]
	if !ok {
		return protocolInfo
	}
	
	layers, ok := layersInterface.([]string)
	if !ok {
		return protocolInfo
	}
	
	// Get highest layer name
	var highestLayer string
	if highestLayerRaw, ok := available["highest_layer"]; ok {
		highestLayer, _ = highestLayerRaw.(string)
	}
	
	// トップレベルプロトコルに基づく処理の振り分け
	// ここではtsharkから直接取得したプロトコル名を優先
	highestLayerLower := strings.ToLower(highestLayer)
	
	// Extract protocol details using appropriate analyzers based on highest layer
	switch highestLayerLower {
	case "http":
		analyzer := protocol_analyzer.NewAnalyzeHTTP(d.Pkt)
		protocolInfo["http_info"] = analyzer.Analyze()
		
	case "dns":
		analyzer := protocol_analyzer.NewAnalyzeDNS(d.Pkt)
		protocolInfo["dns_info"] = analyzer.Analyze()
		
	case "tls", "ssl", "dtls":
		analyzer := protocol_analyzer.NewAnalyzeHTTPS(d.Pkt)
		protocolInfo["https_info"] = analyzer.Analyze()
		
	case "tcp":
		analyzer := protocol_analyzer.NewAnalyzeTCP(d.Pkt)
		protocolInfo["tcp_info"] = analyzer.Analyze()
		
	case "udp":
		analyzer := protocol_analyzer.NewAnalyzeUDP(d.Pkt)
		protocolInfo["udp_info"] = analyzer.Analyze()
		
	case "sctp":
		analyzer := protocol_analyzer.NewAnalyzeSCTP(d.Pkt)
		protocolInfo["sctp_info"] = analyzer.Analyze()
		
	case "arp":
		analyzer := protocol_analyzer.NewAnalyzeARP(d.Pkt)
		protocolInfo["arp_info"] = analyzer.Analyze()
		
	case "icmp", "icmpv6":
		analyzer := protocol_analyzer.NewAnalyzeICMP(d.Pkt)
		protocolInfo["icmp_info"] = analyzer.Analyze()
		
	case "x.25", "x25":
		analyzer := protocol_analyzer.NewAnalyzeX25(d.Pkt)
		protocolInfo["x25_info"] = analyzer.Analyze()
		protocolInfo["protocol_name"] = "X.25"
		
	case "lapb":
		analyzer := protocol_analyzer.NewAnalyzeLAPB(d.Pkt)
		protocolInfo["lapb_info"] = analyzer.Analyze()
		protocolInfo["protocol_name"] = "LAPB"
		
	case "ipars":
		analyzer := protocol_analyzer.NewAnalyzeIPARS(d.Pkt)
		protocolInfo["ipars_info"] = analyzer.Analyze()
		protocolInfo["protocol_name"] = "IPARS"
		
	case "uts":
		analyzer := protocol_analyzer.NewAnalyzeUTS(d.Pkt)
		protocolInfo["uts_info"] = analyzer.Analyze()
		protocolInfo["protocol_name"] = "UTS"
		
	case "q.933", "q933":
		analyzer := protocol_analyzer.NewAnalyzeQ933(d.Pkt)
		protocolInfo["q933_info"] = analyzer.Analyze()
		protocolInfo["protocol_name"] = "Q.933"
		
	default:
		// サポートされていないプロトコルの処理
		// レイヤーベースの検出も合わせて行う
		
		// Extract IP information if present
		if contains(layers, "ip") {
			analyzer := protocol_analyzer.NewAnalyzeIPv4(d.Pkt)
			protocolInfo["ipv4_info"] = analyzer.Analyze()
		}
		
		// サポートされていないプロトコルとして処理
		if len(protocolInfo) <= 1 {  // protocol_nameのみが含まれている場合
			analyzer := protocol_analyzer.NewAnalyzeUnsupportedProtocol(d.Pkt, highestLayer)
			protocolInfo["unsupported_info"] = analyzer.Analyze()
		}
	}
	
	return protocolInfo
}

// Helper function to check if a string is in a slice
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}