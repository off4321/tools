package protocol

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/packet_sequence/models"
)

// IPv4Analyzer はIPv4プロトコルの解析機能を提供する構造体
type IPv4Analyzer struct {
	// デバッグモード
	DebugMode bool
}

// NewIPv4Analyzer は新しいIPv4Analyzerインスタンスを生成する
func NewIPv4Analyzer(debugMode bool) *IPv4Analyzer {
	return &IPv4Analyzer{
		DebugMode: debugMode,
	}
}

// Analyze はIPv4パケットを解析して詳細情報を抽出する
// IPv4の場合、以下の情報を取得：
// - プロトコル番号
// - サービスタイプ
func (a *IPv4Analyzer) Analyze(packet *models.Packet) ([]string, error) {
	if packet.Protocol != "IPv4" {
		return nil, fmt.Errorf("プロトコルがIPv4ではありません: %s", packet.Protocol)
	}

	var details []string
	var protocolNum, tos string

	// 正規表現パターンを定義
	protocolRe := regexp.MustCompile(`Protocol: (\d+) \(([^\)]+)\)`)
	tosRe := regexp.MustCompile(`Differentiated Services Field: (0x[0-9a-fA-F]+)`)

	// パケットの詳細から情報を抽出
	for _, line := range packet.Details {
		// プロトコル番号
		if matches := protocolRe.FindStringSubmatch(line); len(matches) > 2 {
			protocolNum = fmt.Sprintf("%s (%s)", matches[1], matches[2])
		}
		// サービスタイプ
		if matches := tosRe.FindStringSubmatch(line); len(matches) > 1 {
			tos = matches[1]
		}
	}

	// 詳細情報を構築
	if protocolNum != "" {
		details = append(details, fmt.Sprintf("Protocol: %s", protocolNum))
	}
	if tos != "" {
		details = append(details, fmt.Sprintf("ToS: %s", tos))
	}

	// 詳細情報が取得できなかった場合は基本情報のみを返す
	if len(details) == 0 {
		details = append(details, strings.TrimSpace(packet.Info))
	}

	return details, nil
}