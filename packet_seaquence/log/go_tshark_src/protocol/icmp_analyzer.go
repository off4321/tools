package protocol

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/packet_sequence/models"
)

// ICMPAnalyzer はICMPプロトコルの解析機能を提供する構造体
type ICMPAnalyzer struct {
	// デバッグモード
	DebugMode bool
}

// NewICMPAnalyzer は新しいICMPAnalyzerインスタンスを生成する
func NewICMPAnalyzer(debugMode bool) *ICMPAnalyzer {
	return &ICMPAnalyzer{
		DebugMode: debugMode,
	}
}

// Analyze はICMPパケットを解析して詳細情報を抽出する
// ICMPの場合、以下の情報を取得：
// - タイプ (Request/Reply)
func (a *ICMPAnalyzer) Analyze(packet *models.Packet) ([]string, error) {
	if packet.Protocol != "ICMP" {
		return nil, fmt.Errorf("プロトコルがICMPではありません: %s", packet.Protocol)
	}

	var details []string
	var icmpType string
	var typeNumber string

	// 正規表現パターンを定義
	// Type: 8 (Echo (ping) request) または Type: 0 (Echo (ping) reply) の形式を検出
	typeRe := regexp.MustCompile(`Type: (\d+) \(([^\)]+)\)`)
	
	// 詳細情報をクリアして新規作成
	originalDetails := packet.Details
	packet.Details = []string{}

	// パケットの詳細から情報を抽出
	for _, line := range originalDetails {
		// タイプ
		if matches := typeRe.FindStringSubmatch(line); len(matches) > 2 {
			typeNumber = matches[1]
			icmpType = matches[2]
		}
	}

	// タイプ番号でRequestとReplyを判断
	if typeNumber == "8" || strings.Contains(strings.ToLower(icmpType), "request") {
		details = append(details, "Request")
	} else if typeNumber == "0" || strings.Contains(strings.ToLower(icmpType), "reply") {
		details = append(details, "Reply")
	} else if icmpType != "" {
		// その他のタイプはそのまま表示
		details = append(details, fmt.Sprintf("Type: %s", icmpType))
	} else {
		// デフォルトはInfoから決定
		infoLower := strings.ToLower(packet.Info)
		if strings.Contains(infoLower, "echo request") {
			details = append(details, "Request")
		} else if strings.Contains(infoLower, "echo reply") {
			details = append(details, "Reply")
		} else {
			details = append(details, strings.TrimSpace(packet.Info))
		}
	}

	return details, nil
}