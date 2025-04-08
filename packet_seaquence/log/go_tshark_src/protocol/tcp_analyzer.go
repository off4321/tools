package protocol

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/packet_sequence/models"
)

// TCPAnalyzer はTCPプロトコルの解析機能を提供する構造体
type TCPAnalyzer struct {
	// デバッグモード
	DebugMode bool
}

// NewTCPAnalyzer は新しいTCPAnalyzerインスタンスを生成する
func NewTCPAnalyzer(debugMode bool) *TCPAnalyzer {
	return &TCPAnalyzer{
		DebugMode: debugMode,
	}
}

// Analyze はTCPパケットを解析して詳細情報を抽出する
// TCPの場合、以下の情報を取得：
// - ポート番号(Source, Destination)
// - ACK番号
// - コントロールフラグ名
// - ウィンドウサイズ
func (a *TCPAnalyzer) Analyze(packet *models.Packet) ([]string, error) {
	if packet.Protocol != "TCP" {
		return nil, fmt.Errorf("プロトコルがTCPではありません: %s", packet.Protocol)
	}

	var details []string
	var srcPort, dstPort, ack, flags, window string

	// 正規表現パターンを定義
	srcPortRe := regexp.MustCompile(`Source Port: (\d+)`)
	dstPortRe := regexp.MustCompile(`Destination Port: (\d+)`)
	ackRe := regexp.MustCompile(`Acknowledgment number: (\d+)`)
	flagsRe := regexp.MustCompile(`Flags: (0x[0-9a-fA-F]+) \(([^\)]+)\)`)
	windowRe := regexp.MustCompile(`Window size value: (\d+)`)

	// パケットの詳細から情報を抽出
	for _, line := range packet.Details {
		// 送信元ポート
		if matches := srcPortRe.FindStringSubmatch(line); len(matches) > 1 {
			srcPort = matches[1]
		}
		// 宛先ポート
		if matches := dstPortRe.FindStringSubmatch(line); len(matches) > 1 {
			dstPort = matches[1]
		}
		// ACK番号
		if matches := ackRe.FindStringSubmatch(line); len(matches) > 1 {
			ack = matches[1]
		}
		// フラグ
		if matches := flagsRe.FindStringSubmatch(line); len(matches) > 2 {
			flags = matches[2]
		}
		// ウィンドウサイズ
		if matches := windowRe.FindStringSubmatch(line); len(matches) > 1 {
			window = matches[1]
		}
	}

	// 詳細情報を構築
	if srcPort != "" && dstPort != "" {
		details = append(details, fmt.Sprintf("Port: %s→%s", srcPort, dstPort))
	}
	if ack != "" {
		details = append(details, fmt.Sprintf("ACK: %s", ack))
	}
	if flags != "" {
		details = append(details, fmt.Sprintf("Flags: %s", flags))
	}
	if window != "" {
		details = append(details, fmt.Sprintf("Window: %s", window))
	}

	// 詳細情報が取得できなかった場合は基本情報のみを返す
	if len(details) == 0 {
		details = append(details, strings.TrimSpace(packet.Info))
	}

	return details, nil
}