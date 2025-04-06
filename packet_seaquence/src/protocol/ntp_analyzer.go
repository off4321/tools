package protocol

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/packet_sequence/models"
)

// NTPAnalyzer はNTPプロトコルの解析機能を提供する構造体
type NTPAnalyzer struct {
	// デバッグモード
	DebugMode bool
}

// NewNTPAnalyzer は新しいNTPAnalyzerインスタンスを生成する
func NewNTPAnalyzer(debugMode bool) *NTPAnalyzer {
	return &NTPAnalyzer{
		DebugMode: debugMode,
	}
}

// Analyze はNTPパケットを解析して詳細情報を抽出する
// NTPの場合、以下の情報を取得：
// - バージョン
// - モード
// - ストラタム（階層）
func (a *NTPAnalyzer) Analyze(packet *models.Packet) ([]string, error) {
	if packet.Protocol != "NTP" {
		return nil, fmt.Errorf("プロトコルがNTPではありません: %s", packet.Protocol)
	}

	var details []string
	var version, mode, stratum string

	// 正規表現パターンを定義
	versionRe := regexp.MustCompile(`Version: (\d+)`)
	modeRe := regexp.MustCompile(`Mode: ([^\(]+)\((\d+)\)`)
	stratumRe := regexp.MustCompile(`Stratum: (\d+)`)

	// パケットの詳細から情報を抽出
	for _, line := range packet.Details {
		// バージョン
		if matches := versionRe.FindStringSubmatch(line); len(matches) > 1 {
			version = matches[1]
		}
		// モード
		if matches := modeRe.FindStringSubmatch(line); len(matches) > 1 {
			mode = strings.TrimSpace(matches[1])
		}
		// ストラタム
		if matches := stratumRe.FindStringSubmatch(line); len(matches) > 1 {
			stratum = matches[1]
		}
	}

	// 詳細情報を構築
	if version != "" {
		details = append(details, fmt.Sprintf("Version: %s", version))
	}
	if mode != "" {
		details = append(details, fmt.Sprintf("Mode: %s", mode))
	}
	if stratum != "" {
		details = append(details, fmt.Sprintf("Stratum: %s", stratum))
	}

	// 詳細情報が取得できなかった場合は基本情報のみを返す
	if len(details) == 0 {
		details = append(details, strings.TrimSpace(packet.Info))
	}

	return details, nil
}