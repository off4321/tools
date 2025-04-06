package protocol

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/packet_sequence/models"
)

// HTTPAnalyzer はHTTPプロトコルの解析機能を提供する構造体
type HTTPAnalyzer struct {
	// デバッグモード
	DebugMode bool
}

// NewHTTPAnalyzer は新しいHTTPAnalyzerインスタンスを生成する
func NewHTTPAnalyzer(debugMode bool) *HTTPAnalyzer {
	return &HTTPAnalyzer{
		DebugMode: debugMode,
	}
}

// Analyze はHTTPパケットを解析して詳細情報を抽出する
// HTTPの場合、以下の情報を取得：
// - HTTPメソッド
// - ステータスコード
// - ステータスメッセージ
// - バージョン
func (a *HTTPAnalyzer) Analyze(packet *models.Packet) ([]string, error) {
	if packet.Protocol != "HTTP" {
		return nil, fmt.Errorf("プロトコルがHTTPではありません: %s", packet.Protocol)
	}

	var details []string
	var method, statusCode, statusMsg, version string

	// 正規表現パターンを定義
	requestLineRe := regexp.MustCompile(`(GET|POST|PUT|DELETE|HEAD|OPTIONS|TRACE|CONNECT|PATCH) ([^ ]+) (HTTP/[0-9.]+)`)
	statusLineRe := regexp.MustCompile(`(HTTP/[0-9.]+) (\d+) ([^\r\n]+)`)

	// パケットの詳細から情報を抽出
	for _, line := range packet.Details {
		// リクエスト行
		if matches := requestLineRe.FindStringSubmatch(line); len(matches) > 3 {
			method = matches[1]
			 // URIは変数を宣言せずに直接使用
			version = matches[3]
		}
		// ステータス行
		if matches := statusLineRe.FindStringSubmatch(line); len(matches) > 3 {
			version = matches[1]
			statusCode = matches[2]
			statusMsg = matches[3]
		}
	}

	// 詳細情報を構築
	if method != "" {
		details = append(details, fmt.Sprintf("Method: %s", method))
	}
	
	// URLは表示しない (長いため)
	
	if statusCode != "" {
		if statusMsg != "" {
			details = append(details, fmt.Sprintf("Status: %s %s", statusCode, statusMsg))
		} else {
			details = append(details, fmt.Sprintf("Status: %s", statusCode))
		}
	}
	if version != "" && len(details) == 0 {
		details = append(details, fmt.Sprintf("Version: %s", version))
	}

	// 詳細情報が取得できなかった場合は基本情報のみを返す
	if len(details) == 0 {
		details = append(details, strings.TrimSpace(packet.Info))
	}

	return details, nil
}