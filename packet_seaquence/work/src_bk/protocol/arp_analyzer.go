package protocol

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/packet_sequence/models"
)

// ARPAnalyzer はARPプロトコルの解析機能を提供する構造体
type ARPAnalyzer struct {
	// デバッグモード
	DebugMode bool
}

// NewARPAnalyzer は新しいARPAnalyzerインスタンスを生成する
func NewARPAnalyzer(debugMode bool) *ARPAnalyzer {
	return &ARPAnalyzer{
		DebugMode: debugMode,
	}
}

// Analyze はARPパケットを解析して詳細情報を抽出する
// ARPの場合、以下の情報を取得：
// - オペレーションコード
// - 送信元IPアドレスとMACアドレス
// - 目標IPアドレスとMACアドレス
func (a *ARPAnalyzer) Analyze(packet *models.Packet) ([]string, error) {
	if packet.Protocol != "ARP" {
		return nil, fmt.Errorf("プロトコルがARPではありません: %s", packet.Protocol)
	}

	var details []string
	var opcode, senderIP, targetIP, targetMAC string

	// 正規表現パターンを定義
	opcodeRe := regexp.MustCompile(`Opcode: ([^\(]+)`)
	senderIPRe := regexp.MustCompile(`Sender IP address: ([0-9.]+)`)
	targetIPRe := regexp.MustCompile(`Target IP address: ([0-9.]+)`)
	targetMACRe := regexp.MustCompile(`Target MAC address: ([0-9A-Fa-f:]+)`)

	// パケットの詳細から情報を抽出
	for _, line := range packet.Details {
		// オペレーションコード
		if matches := opcodeRe.FindStringSubmatch(line); len(matches) > 1 {
			opcode = strings.TrimSpace(matches[1])
		}
		// 送信元IPアドレス
		if matches := senderIPRe.FindStringSubmatch(line); len(matches) > 1 {
			senderIP = matches[1]
		}
		// 目標IPアドレス
		if matches := targetIPRe.FindStringSubmatch(line); len(matches) > 1 {
			targetIP = matches[1]
		}
		// 目標MACアドレス
		if matches := targetMACRe.FindStringSubmatch(line); len(matches) > 1 {
			targetMAC = matches[1]
		}
	}

	// ARPパケットでは常に送信元と宛先IPアドレスを更新
	if senderIP != "" {
		packet.Source = senderIP
	}
	if targetIP != "" {
		packet.Destination = targetIP
	}

	// 詳細情報を構築
	if opcode != "" {
		details = append(details, fmt.Sprintf("Opcode: %s", opcode))
	}

	if targetMAC != "" {
		// 00:00:00:00:00:00はブロードキャストまたは未知のアドレスを示す
		if targetMAC == "00:00:00:00:00:00" {
			// ARPリクエストの場合は特別なメッセージを表示
			if strings.Contains(opcode, "request") {
				details = append(details, fmt.Sprintf("Target MAC: Broadcast (seeking %s)", targetIP))
			} else {
				details = append(details, "Target MAC: Unknown")
			}
		} else {
			details = append(details, fmt.Sprintf("Target MAC: %s", targetMAC))
		}
	}

	// 詳細情報が取得できなかった場合は基本情報のみを返す
	if len(details) == 0 {
		details = append(details, strings.TrimSpace(packet.Info))
	}

	return details, nil
}
