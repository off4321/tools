package protocol

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/packet_sequence/models"
)

// DNSAnalyzer はDNSプロトコルの解析機能を提供する構造体
type DNSAnalyzer struct {
	// デバッグモード
	DebugMode bool
}

// NewDNSAnalyzer は新しいDNSAnalyzerインスタンスを生成する
func NewDNSAnalyzer(debugMode bool) *DNSAnalyzer {
	return &DNSAnalyzer{
		DebugMode: debugMode,
	}
}

// Analyze はDNSパケットを解析して詳細情報を抽出する
// DNSの場合、以下の情報を取得：
// - クエリ名
// - クエリタイプ
// - 応答数
// - リプライコード
func (a *DNSAnalyzer) Analyze(packet *models.Packet) ([]string, error) {
	if packet.Protocol != "DNS" {
		return nil, fmt.Errorf("プロトコルがDNSではありません: %s", packet.Protocol)
	}

	var details []string
	var queryName, queryType, replyCode string

	// レスポンスかどうかを判定
	isResponse := strings.Contains(packet.Info, "response")

	// レスポンスのカウント用
	aRecordCount := 0
	nsRecordCount := 0
	aaaaRecordCount := 0
	cnameRecordCount := 0
	otherRecordCount := 0

	// 正規表現パターンを定義
	queryNameRe := regexp.MustCompile(`Name: ([^\s]+)`)
	queryTypeRe := regexp.MustCompile(`Type: ([A-Z]+) \(`)
	replyCodeRe := regexp.MustCompile(`Reply code: ([^(]+)`)
	altReplyCodeRe := regexp.MustCompile(`= Reply code: ([^(]+)`)
	aRecordRe := regexp.MustCompile(`Type: A \(`)
	nsRecordRe := regexp.MustCompile(`Type: NS \(`)
	aaaaRecordRe := regexp.MustCompile(`Type: AAAA \(`)
	cnameRecordRe := regexp.MustCompile(`Type: CNAME \(`)

	// パケットの詳細から情報を抽出
	queriesFound := false
	answersFound := false

	for _, line := range packet.Details {

		// リプライコード (2種類のパターンに対応)
		if matches := replyCodeRe.FindStringSubmatch(line); len(matches) > 1 {
			replyCode = strings.TrimSpace(matches[1])
		}
		if replyCode == "" {
			if matches := altReplyCodeRe.FindStringSubmatch(line); len(matches) > 1 {
				replyCode = strings.TrimSpace(matches[1])
			}
		}

		// クエリセクションを検出
		if strings.Contains(line, "Queries") {
			queriesFound = true
			continue
		}

		// 応答セクションを検出
		if strings.Contains(line, "Answers") {
			queriesFound = false
			answersFound = true
			continue
		}

		// クエリ名とタイプを抽出 (クエリセクション内または単独の場合)
		// 名前
		if matches := queryNameRe.FindStringSubmatch(line); len(matches) > 1 {
			if queryName == "" || queriesFound { // まだ設定されていないか、クエリセクション内なら更新
				queryName = matches[1]
			}
		}
		// タイプ
		if matches := queryTypeRe.FindStringSubmatch(line); len(matches) > 1 {
			// クエリセクション内か、まだタイプが設定されていない場合に更新
			if queriesFound || (queryType == "" && !answersFound) {
				queryType = matches[1]
			}
		}

		// レスポンスの場合、レコードタイプをカウント
		if isResponse && answersFound {
			if aRecordRe.MatchString(line) {
				aRecordCount++
			} else if nsRecordRe.MatchString(line) {
				nsRecordCount++
			} else if aaaaRecordRe.MatchString(line) {
				aaaaRecordCount++
			} else if cnameRecordRe.MatchString(line) {
				cnameRecordCount++
			} else if queryTypeRe.MatchString(line) {
				// その他のレコードタイプ
				otherRecordCount++
			}
		}
	}

	// 詳細情報を構築
	if isResponse {
		// レスポンスパケットの場合
		if queryName != "" && queryType != "" {
			details = append(details, "Response - Query: "+queryName+", Type: "+queryType)
		} else {
			details = append(details, "Response")
		}

		if replyCode != "" {
			details = append(details, "Reply code: "+replyCode)
		}
	} else {
		// クエリパケットの場合
		if queryName != "" && queryType != "" {
			details = append(details, "Query - "+queryName+", Type: "+queryType)
		} else {
			// Info欄から情報を抽出する試み
			infoRe := regexp.MustCompile(`Standard query .+ ([A-Z]+) ([^ ]+)`)
			if matches := infoRe.FindStringSubmatch(packet.Info); len(matches) > 2 {
				queryType = matches[1]
				queryName = matches[2]
				details = append(details, "Query - "+queryName+", Type: "+queryType)
			}
		}
	}

	// 詳細情報が取得できなかった場合は基本情報のみを返す
	if len(details) == 0 {
		details = append(details, strings.TrimSpace(packet.Info))
	}

	return details, nil
}
