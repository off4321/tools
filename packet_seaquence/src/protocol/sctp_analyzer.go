package protocol

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/packet_sequence/models"
)

// SCTPAnalyzer はSCTPプロトコルの解析機能を提供する構造体
type SCTPAnalyzer struct {
	// デバッグモード
	DebugMode bool
}

// NewSCTPAnalyzer は新しいSCTPAnalyzerインスタンスを生成する
func NewSCTPAnalyzer(debugMode bool) *SCTPAnalyzer {
	return &SCTPAnalyzer{
		DebugMode: debugMode,
	}
}

// Analyze はSCTPパケットを解析して詳細情報を抽出する
// SCTPの場合、以下の情報を取得：
// - ポート番号(Source, Destination)
// - チャンクタイプ
func (a *SCTPAnalyzer) Analyze(packet *models.Packet) ([]string, error) {
	if packet.Protocol != "SCTP" {
		return nil, fmt.Errorf("プロトコルがSCTPではありません: %s", packet.Protocol)
	}

	var details []string
	var srcPort, dstPort, chunkType string

	// 正規表現パターンを定義
	srcPortRe := regexp.MustCompile(`Source Port: (\d+)`)
	dstPortRe := regexp.MustCompile(`Destination Port: (\d+)`)
	chunkTypeRe := regexp.MustCompile(`Chunk Type: ([^\(]+)\((\d+)\)`)

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
		// チャンクタイプ
		if matches := chunkTypeRe.FindStringSubmatch(line); len(matches) > 1 {
			chunkType = strings.TrimSpace(matches[1])
		}
	}

	// 詳細情報を構築
	if srcPort != "" && dstPort != "" {
		details = append(details, fmt.Sprintf("Port: %s→%s", srcPort, dstPort))
	}
	if chunkType != "" {
		details = append(details, fmt.Sprintf("Chunk: %s", chunkType))
	}

	// 詳細情報が取得できなかった場合は基本情報のみを返す
	if len(details) == 0 {
		details = append(details, strings.TrimSpace(packet.Info))
	}

	return details, nil
}