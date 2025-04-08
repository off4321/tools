// Package reader provides PCAPからのパケット読み込み機能
package reader

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"
)

// PCAPReader はPCAPファイルの読み込みを行う
type PCAPReader struct {
	Filename string
}

// NewPCAPReader creates a new PCAPReader
func NewPCAPReader(filename string) *PCAPReader {
	return &PCAPReader{
		Filename: filename,
	}
}

// Open opens the PCAP file (no-op for tshark implementation as we open on read)
func (r *PCAPReader) Open() error {
	// tsharkの存在を確認
	_, err := exec.LookPath("tshark")
	if err != nil {
		return fmt.Errorf("tsharkが見つかりません。インストールしてください: %v", err)
	}
	
	// PCAPファイルの存在確認はtshark実行時に自動的に行われる
	fmt.Printf("PCAPファイル読み込み準備完了: %s\n", r.Filename)
	return nil
}

// Close closes the PCAP file (no-op for tshark implementation)
func (r *PCAPReader) Close() error {
	return nil
}

// Packet represents a network packet with its layers and metadata
type Packet struct {
	Layers       []string
	HighestLayer string
	Timestamp    time.Time
	SrcIP        string
	DstIP        string
	Protocol     string
	FrameLen     int
	Info         map[string]interface{}
	RawInfo      string
}

// LayerExists checks if a layer exists in the packet
func (p *Packet) LayerExists(layerName string) bool {
	layerName = strings.ToLower(layerName)
	for _, l := range p.Layers {
		if strings.ToLower(l) == layerName {
			return true
		}
	}
	return false
}

// GetField gets a field from the packet info
func (p *Packet) GetField(layer, field string) (interface{}, bool) {
	if p.Info == nil {
		return nil, false
	}
	layerInfo, ok := p.Info[layer]
	if !ok {
		return nil, false
	}
	layerMap, ok := layerInfo.(map[string]interface{})
	if !ok {
		return nil, false
	}
	fieldValue, ok := layerMap[field]
	return fieldValue, ok
}

// SetField sets a field in the packet info
func (p *Packet) SetField(layer, field string, value interface{}) {
	if p.Info == nil {
		p.Info = make(map[string]interface{})
	}
	layerInfo, ok := p.Info[layer]
	if !ok {
		p.Info[layer] = map[string]interface{}{field: value}
		return
	}
	layerMap, ok := layerInfo.(map[string]interface{})
	if !ok {
		// 既存のレイヤー情報が想定されるマップではない場合、新しいマップとして作り直す
		p.Info[layer] = map[string]interface{}{field: value}
		return
	}
	layerMap[field] = value
	p.Info[layer] = layerMap
}

// AddLayer adds a layer to the packet
func (p *Packet) AddLayer(layerName string) {
	p.Layers = append(p.Layers, layerName)
}

// GetProtocolName returns the protocol name of the packet
func (p *Packet) GetProtocolName() string {
	// First try to get from highest layer or protocol field
	if p.Protocol != "" {
		return p.Protocol
	}
	
	if p.HighestLayer != "" && p.HighestLayer != "Payload" && p.HighestLayer != "DecodeFailure" {
		return p.HighestLayer
	}
	
	// Try to get from custom protocols
	if p.LayerExists("x25") {
		return "X25"
	} else if p.LayerExists("lapb") {
		return "LAPB"
	} else if p.LayerExists("ipars") {
		return "IPARS"
	}
	
	// Default fallback
	if len(p.Layers) > 0 {
		return p.Layers[len(p.Layers)-1]
	}
	
	return "Unknown"
}

// GetTimestamp returns the timestamp of the packet
func (p *Packet) GetTimestamp() time.Time {
	return p.Timestamp
}

// Read reads packets from the PCAP file using tshark
func (r *PCAPReader) Read(maxPackets int, verbose bool) ([]Packet, error) {
	// tsharkコマンドを構築
	args := []string{
		"-r", r.Filename,
		"-T", "fields",
		"-E", "separator=|",
		"-E", "header=y",
		"-e", "frame.time_epoch",
		"-e", "frame.len",
		"-e", "ip.src",
		"-e", "ip.dst",
		"-e", "frame.protocols",
		"-e", "_ws.col.Protocol",
		"-e", "_ws.col.Info",
	}

	if verbose {
		fmt.Printf("実行するtsharkコマンド: tshark %s\n", strings.Join(args, " "))
	}

	// コマンドを実行
	cmd := exec.Command("tshark", args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("tshark出力パイプの作成に失敗しました: %v", err)
	}

	err = cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("tsharkの実行に失敗しました: %v", err)
	}

	// 出力を解析してパケットのスライスを作成
	packets := []Packet{}
	reader := bufio.NewReader(stdout)
	
	// ヘッダー行を読み取り
	headerLine, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("tsharkの出力ヘッダーの読み取りに失敗しました: %v", err)
	}
	
	// ヘッダーフィールドの解析
	headers := strings.Split(strings.TrimSpace(headerLine), "|")
	if verbose {
		fmt.Printf("tsharkヘッダー: %v\n", headers)
	}
	
	count := 0
	for {
		if maxPackets > 0 && count >= maxPackets {
			break
		}

		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("tsharkの出力の読み取りに失敗しました: %v", err)
		}

		// 改行文字を削除
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// パイプで区切られたフィールドを解析
		fields := strings.Split(line, "|")
		if len(fields) < 7 {
			if verbose {
				fmt.Printf("警告: フィールド数が不足しているパケット行をスキップします: %s\n", line)
			}
			continue
		}

		// タイムスタンプをパース
		timestampFloat := 0.0
		fmt.Sscanf(fields[0], "%f", &timestampFloat)
		timestamp := time.Unix(int64(timestampFloat), int64((timestampFloat-float64(int64(timestampFloat)))*1e9))

		// フレーム長をパース
		frameLen := 0
		fmt.Sscanf(fields[1], "%d", &frameLen)

		// 送信元IPと宛先IP
		srcIP := fields[2]
		dstIP := fields[3]

		// プロトコルスタックとトッププロトコル
		protocolStack := strings.Split(fields[4], ":")
		topProtocol := fields[5]
		
		// パケット情報
		packetInfo := fields[6]

		// Packetオブジェクトを作成
		packet := Packet{
			Timestamp:    timestamp,
			SrcIP:        srcIP,
			DstIP:        dstIP,
			FrameLen:     frameLen,
			Layers:       protocolStack,
			HighestLayer: topProtocol,
			Protocol:     topProtocol,
			Info:         make(map[string]interface{}),
			RawInfo:      packetInfo,
		}

		// プロトコルスタックからレイヤー情報を設定
		// Ensure protocol stack is populated with values tshark returns
		for _, proto := range protocolStack {
			proto = strings.TrimSpace(proto)
			if proto != "" {
				packet.AddLayer(proto)
			}
		}

		// プロトコル情報を設定（トッププロトコルを基準に）
		topProtocolLower := strings.ToLower(topProtocol)
		
		// 基本情報の設定
		packet.SetField("general", "frame_len", frameLen)
		packet.SetField("general", "timestamp", timestamp.String())
		packet.SetField("general", "info", packetInfo)
		
		// IPレイヤー情報の設定
		if srcIP != "" && dstIP != "" {
			packet.SetField("ip", "src", srcIP)
			packet.SetField("ip", "dst", dstIP)
			
			// ipv4_infoも設定（互換性のため）
			ipv4Info := map[string]interface{}{
				"src": srcIP,
				"dst": dstIP,
			}
			packet.Info["ipv4_info"] = ipv4Info
		}
		
		// サポートされているプロトコルに応じて特定のフィールドを設定
		// これにより、以前の処理との互換性を確保
		switch topProtocolLower {
		case "http":
			packet.SetField("http", "request_or_response", packetInfo)
			packet.Info["http_info"] = map[string]interface{}{
				"info": packetInfo,
			}
		case "dns":
			packet.SetField("dns", "query_or_response", packetInfo)
			packet.Info["dns_info"] = map[string]interface{}{
				"info": packetInfo,
			}
		case "tcp":
			packet.SetField("tcp", "info", packetInfo)
			
			// 送信元・宛先ポートの抽出（例: "80 → 36000"）
			if ports := extractPorts(packetInfo); len(ports) == 2 {
				packet.SetField("tcp", "srcport", ports[0])
				packet.SetField("tcp", "dstport", ports[1])
			}
			
			packet.Info["tcp_info"] = map[string]interface{}{
				"info": packetInfo,
			}
		case "udp":
			packet.SetField("udp", "info", packetInfo)
			
			// 送信元・宛先ポートの抽出
			if ports := extractPorts(packetInfo); len(ports) == 2 {
				packet.SetField("udp", "srcport", ports[0])
				packet.SetField("udp", "dstport", ports[1])
			}
			
			packet.Info["udp_info"] = map[string]interface{}{
				"info": packetInfo,
			}
		case "sctp":
			packet.SetField("sctp", "info", packetInfo)
			packet.Info["sctp_info"] = map[string]interface{}{
				"info": packetInfo,
			}
		case "arp":
			packet.SetField("arp", "info", packetInfo)
			packet.Info["arp_info"] = map[string]interface{}{
				"info": packetInfo,
			}
			
			// ARPの送信元・宛先情報の抽出（IPが存在しない場合の対応）
			if srcIP == "" || dstIP == "" {
				arpSrc, arpDst := extractArpAddresses(packetInfo)
				if arpSrc != "" {
					packet.SrcIP = arpSrc
					packet.SetField("arp", "src_proto_ipv4", arpSrc)
				}
				if arpDst != "" {
					packet.DstIP = arpDst
					packet.SetField("arp", "dst_proto_ipv4", arpDst)
				}
			}
		case "icmp":
			packet.SetField("icmp", "info", packetInfo)
			packet.Info["icmp_info"] = map[string]interface{}{
				"info": packetInfo,
				"type_name": extractIcmpType(packetInfo),
			}
		default:
			// その他のプロトコル（サポートされていないもの）
			packet.SetField("unsupported", "protocol_name", topProtocol)
			packet.SetField("unsupported", "info", packetInfo)
			
			// src, dstが設定されていない場合の対応
			if srcIP == "" {
				packet.SrcIP = "Unknown"
			}
			if dstIP == "" {
				packet.DstIP = "Remote"
			}
			
			// 追加情報の抽出
			if textInfo := extractTextInfo(packetInfo); textInfo != "" {
				packet.SetField("unsupported", "message", textInfo)
			}
		}

		packets = append(packets, packet)
		count++
	}

	err = cmd.Wait()
	if err != nil {
		return nil, fmt.Errorf("tsharkコマンドの実行に失敗しました: %v", err)
	}

	fmt.Printf("読み込み完了: %dパケット\n", count)
	return packets, nil
}

// extractPorts extracts source and destination ports from packet info string
func extractPorts(info string) []string {
	// 一般的なport表記: "80 → 36000" or "80→36000"
	portParts := strings.Split(strings.Replace(info, " ", "", -1), "→")
	if len(portParts) == 2 {
		// ポート部分だけを抽出（他の情報が含まれている可能性がある）
		srcPort := extractNumber(portParts[0])
		dstPort := extractNumber(portParts[1])
		if srcPort != "" && dstPort != "" {
			return []string{srcPort, dstPort}
		}
	}
	return []string{}
}

// extractNumber extracts the first number from a string
func extractNumber(s string) string {
	var number string
	inNumber := false
	
	for _, c := range s {
		if c >= '0' && c <= '9' {
			number += string(c)
			inNumber = true
		} else if inNumber {
			// 数字の後に数字以外が来たら終了
			break
		}
	}
	
	return number
}

// extractArpAddresses extracts source and destination addresses from ARP info
func extractArpAddresses(info string) (string, string) {
	// 一般的なARP表記: "Who has 192.168.1.1? Tell 192.168.1.100"
	// または "192.168.1.100 is at 00:11:22:33:44:55"
	var src, dst string
	
	if strings.Contains(info, "Who has") && strings.Contains(info, "Tell") {
		parts := strings.Split(info, "Tell")
		if len(parts) == 2 {
			// "Tell"の後のIPアドレスが送信元
			src = extractIpAddress(parts[1])
			
			// "Who has"の後のIPアドレスが宛先
			whoParts := strings.Split(parts[0], "Who has")
			if len(whoParts) == 2 {
				dst = extractIpAddress(whoParts[1])
			}
		}
	} else if strings.Contains(info, "is at") {
		parts := strings.Split(info, "is at")
		if len(parts) == 2 {
			// "is at"の前のIPアドレスが送信元
			src = extractIpAddress(parts[0])
			// ARPリプライのdstはMACアドレスで表記されるため、空とする
			dst = ""
		}
	}
	
	return src, dst
}

// extractIpAddress extracts the first IP address from a string
func extractIpAddress(s string) string {
	// IPアドレスの正規表現による抽出の簡易実装
	parts := strings.Fields(s)
	for _, part := range parts {
		// ドットを含む文字列を検索（IPアドレスの簡易判定）
		if strings.Count(part, ".") == 3 {
			// ピリオドで終わる場合は除去
			if strings.HasSuffix(part, ".") {
				part = part[:len(part)-1]
			}
			// カンマで終わる場合は除去
			if strings.HasSuffix(part, ",") {
				part = part[:len(part)-1]
			}
			// 疑問符で終わる場合は除去
			if strings.HasSuffix(part, "?") {
				part = part[:len(part)-1]
			}
			return part
		}
	}
	return ""
}

// extractIcmpType extracts ICMP type name from packet info
func extractIcmpType(info string) string {
	// 一般的なICMP type表記: "Echo (ping) request", "Echo (ping) reply"
	if strings.Contains(info, "Echo") {
		if strings.Contains(info, "request") {
			return "Echo Request"
		} else if strings.Contains(info, "reply") {
			return "Echo Reply"
		}
	}
	// その他のICMPタイプの抽出（必要に応じて追加）
	return "Unknown"
}

// extractTextInfo extracts textual information from packet info
func extractTextInfo(info string) string {
	// 特定のプロトコル情報から有用なテキストを抽出
	// 情報自体がすでにテキストとなっている場合は、そのまま返す
	if len(info) > 0 {
		return info
	}
	return ""
}