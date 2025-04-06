package protocol

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/packet_sequence/models"
)

// UDPAnalyzer はUDPプロトコルの解析機能を提供する構造体
type UDPAnalyzer struct {
	// デバッグモード
	DebugMode bool
}

// NewUDPAnalyzer は新しいUDPAnalyzerインスタンスを生成する
func NewUDPAnalyzer(debugMode bool) *UDPAnalyzer {
	return &UDPAnalyzer{
		DebugMode: debugMode,
	}
}

// Analyze はUDPパケットを解析して詳細情報を抽出する
// UDPの場合、以下の情報を取得：
// - ポート番号(Source, Destination)
func (a *UDPAnalyzer) Analyze(packet *models.Packet) ([]string, error) {
	if packet.Protocol != "UDP" {
		return nil, fmt.Errorf("プロトコルがUDPではありません: %s", packet.Protocol)
	}

	var details []string
	var srcPort, dstPort string

	// 正規表現パターンを定義
	srcPortRe := regexp.MustCompile(`Source Port: (\d+)`)
	dstPortRe := regexp.MustCompile(`Destination Port: (\d+)`)

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
	}

	// 詳細情報を構築
	if srcPort != "" && dstPort != "" {
		details = append(details, fmt.Sprintf("Port: %s→%s", srcPort, dstPort))
	}

	// 詳細情報が取得できなかった場合は基本情報のみを返す
	if len(details) == 0 {
		details = append(details, strings.TrimSpace(packet.Info))
	}

	return details, nil
}