package protocol

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/packet_sequence/models"
)

// X25Analyzer はX.25プロトコルの解析機能を提供する構造体
type X25Analyzer struct {
	// デバッグモード
	DebugMode bool
}

// NewX25Analyzer は新しいX25Analyzerインスタンスを生成する
func NewX25Analyzer(debugMode bool) *X25Analyzer {
	return &X25Analyzer{
		DebugMode: debugMode,
	}
}

// Analyze はX.25パケットを解析して詳細情報を抽出する
// X.25の場合、以下の情報を取得：
// - パケット種別
// - チャネル番号
// - DLCI
func (a *X25Analyzer) Analyze(packet *models.Packet) ([]string, error) {
	if packet.Protocol != "X.25" {
		return nil, fmt.Errorf("プロトコルがX.25ではありません: %s", packet.Protocol)
	}

	var details []string
	var channelNum, dlci string

	// 正規表現パターンを定義
	packetTypeRe := regexp.MustCompile(`\b(CN|CR|CC|CA|CI|CQ|CF|DT|IT|IF|RR|RNR|REJ|RI|RQ|RF|SI|SQ|SF|Diagnostic|Registration Confirmation|Registration Request)\b`)
	dataRe := regexp.MustCompile(`Data VC:(\d+) P\(S\):(\d+) P\(R\):(\d+)`)
	packetTypeDataRe := regexp.MustCompile(`Packet Type: Data \(0x0\)`)
	logicalChannelRe := regexp.MustCompile(`Logical Channel: (\d+)`)
	prRe := regexp.MustCompile(`P\(R\): (\d+)`)
	psRe := regexp.MustCompile(`P\(S\): (\d+)`)
	channelNumRe := regexp.MustCompile(`Channel: (\d+)`)
	dlciRe := regexp.MustCompile(`DLCI: (\d+)`)
	vcRe := regexp.MustCompile(`VC:(\d+)`)
	funcRe := regexp.MustCompile(`func=([A-Z]+)`)

	// Data パケットフラグ（重複出力を防ぐため）
	isDataPacket := false

	// パケットの詳細から情報を抽出（検出のみ）
	var packetTypeValue string
	var vcValue string
	var psValue string
	var prValue string
	var channelValue string
	var dlciValue string
	var funcValue string

	for _, line := range packet.Details {
		// データパケット検出（tsharkの詳細出力形式）
		if packetTypeDataRe.MatchString(line) {
			isDataPacket = true
			packetTypeValue = "Data"
		}
		
		// 論理チャネル番号
		if matches := logicalChannelRe.FindStringSubmatch(line); len(matches) > 1 {
			vcValue = matches[1]
		}
		
		// P(R)値
		if matches := prRe.FindStringSubmatch(line); len(matches) > 1 && isDataPacket {
			prValue = matches[1]
		}
		
		// P(S)値
		if matches := psRe.FindStringSubmatch(line); len(matches) > 1 && isDataPacket {
			psValue = matches[1]
		}
		
		// Data パケット（以前の形式）
		if matches := dataRe.FindStringSubmatch(line); len(matches) > 1 {
			packetTypeValue = "Data"
			vcValue = matches[1]
			psValue = matches[2]
			prValue = matches[3]
			isDataPacket = true
		}
		
		// 他のパケット種別 (例: CR, RR など)
		if !isDataPacket {
			if typeMatches := packetTypeRe.FindStringSubmatch(line); len(typeMatches) > 1 {
				packetTypeValue = typeMatches[1]
			}
		}
		
		// チャネル番号
		if matches := channelNumRe.FindStringSubmatch(line); len(matches) > 1 {
			channelValue = matches[1]
		}
		
		// DLCI
		if matches := dlciRe.FindStringSubmatch(line); len(matches) > 1 {
			dlciValue = matches[1]
		}
		
		// パケット種別 (例: RR)
		if matches := funcRe.FindStringSubmatch(line); len(matches) > 1 {
			funcValue = matches[1]
		}
	}

	// X.25プロトコルでDataパケットを検出した場合
	if isDataPacket || (packet.Protocol == "X.25" && strings.Contains(packet.Info, "Data")) {
		packetTypeValue = "Data"
	}

	// VC番号 (重複を避けるため、Data パケットの場合は既に追加済み)
	if vcValue == "" && !isDataPacket {
		for _, line := range packet.Details {
			if matches := vcRe.FindStringSubmatch(line); len(matches) > 1 {
				vcValue = matches[1]
				break
			}
		}
	}

	// 結果を一貫した順序で構築
	if packetTypeValue != "" {
		details = append(details, fmt.Sprintf("Packet Type: %s", packetTypeValue))
	}
	
	if vcValue != "" {
		details = append(details, fmt.Sprintf("VC: %s", vcValue))
	}
	
	if channelValue != "" || channelNum != "" {
		value := channelValue
		if value == "" {
			value = channelNum
		}
		details = append(details, fmt.Sprintf("Channel: %s", value))
	}
	
	if psValue != "" {
		details = append(details, fmt.Sprintf("P(S): %s", psValue))
	}
	
	if prValue != "" {
		details = append(details, fmt.Sprintf("P(R): %s", prValue))
	}
	
	if dlciValue != "" || dlci != "" {
		value := dlciValue
		if value == "" {
			value = dlci
		}
		details = append(details, fmt.Sprintf("DLCI: %s", value))
	}
	
	if funcValue != "" {
		details = append(details, fmt.Sprintf("Func: %s", funcValue))
	}

	// 詳細情報が取得できなかった場合は基本情報のみを返す
	if len(details) == 0 {
		details = append(details, strings.TrimSpace(packet.Info))
	}

	return details, nil
}