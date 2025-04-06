package writer

import (
	"fmt"
	"net"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/packet_sequence/models"
)

// Writer はパケット情報を出力するインターフェース
type Writer interface {
	Write(packets []*models.Packet) error
}

// MermaidWriter はMermaid形式でパケット情報をMarkdownファイルに出力する構造体
type MermaidWriter struct {
	OutputPath string // 出力ファイルのパス
	DebugMode  bool   // デバッグモード
}

// NewMermaidWriter は新しいMermaidWriterを生成する
func NewMermaidWriter(outputPath string, debugMode bool) *MermaidWriter {
	return &MermaidWriter{
		OutputPath: outputPath,
		DebugMode:  debugMode,
	}
}

// IPアドレスを比較する関数（ソート用）
func compareIP(ip1, ip2 string) bool {
	// IPアドレスをパースしてバイト列に変換
	netIP1 := net.ParseIP(ip1)
	netIP2 := net.ParseIP(ip2)

	// どちらもIPアドレスでない場合は通常の文字列比較
	if netIP1 == nil && netIP2 == nil {
		return ip1 < ip2
	}

	// IP1がIPアドレスでない場合はIP2の方が大きいとみなす
	if netIP1 == nil {
		return true
	}

	// IP2がIPアドレスでない場合はIP1の方が大きいとみなす
	if netIP2 == nil {
		return false
	}

	// どちらもIPアドレスの場合はバイト列を比較
	return compareIPBytes(netIP1, netIP2)
}

// IPアドレスのバイト列を比較
func compareIPBytes(ip1, ip2 net.IP) bool {
	// IPv4アドレス同士の比較
	ip1v4 := ip1.To4()
	ip2v4 := ip2.To4()

	// どちらもIPv4の場合
	if ip1v4 != nil && ip2v4 != nil {
		for i := 0; i < 4; i++ {
			if ip1v4[i] != ip2v4[i] {
				return ip1v4[i] < ip2v4[i]
			}
		}
		return false // 同じIPアドレス
	}

	// IPv4とIPv6の場合はIPv4を小さいとみなす
	if ip1v4 != nil && ip2v4 == nil {
		return true
	}

	if ip1v4 == nil && ip2v4 != nil {
		return false
	}

	// どちらもIPv6の場合
	for i := 0; i < 16; i++ {
		if ip1[i] != ip2[i] {
			return ip1[i] < ip2[i]
		}
	}
	return false // 同じIPアドレス
}

// Write はパケット情報をMermaid形式でMarkdownファイルに出力する
func (w *MermaidWriter) Write(packets []*models.Packet) error {
	// ファイルを開く
	file, err := os.Create(w.OutputPath)
	if err != nil {
		return fmt.Errorf("ファイル作成エラー: %v", err)
	}
	defer file.Close()

	// Markdownヘッダーを書き込む
	header := "# パケットシーケンス図\n\n"
	file.WriteString(header)

	// Mermaid形式でシーケンス図を書き込む
	mermaidHeader := "```mermaid\nsequenceDiagram\n"
	file.WriteString(mermaidHeader)

	// パケットからユニークなIPアドレスを収集
	actorSet := make(map[string]bool)
	var actors []string
	for _, packet := range packets {
		if !actorSet[packet.Source] {
			actors = append(actors, packet.Source)
			actorSet[packet.Source] = true
		}
		if !actorSet[packet.Destination] {
			actors = append(actors, packet.Destination)
			actorSet[packet.Destination] = true
		}
	}

	// IPアドレスを昇順にソート
	sort.Slice(actors, func(i, j int) bool {
		return compareIP(actors[i], actors[j])
	})

	// ソートされたアクターを追加
	for _, actor := range actors {
		file.WriteString(fmt.Sprintf("    participant %s\n", actor))
	}

	// パケットごとにシーケンス図のメッセージを追加
	for _, packet := range packets {
		// メッセージラベルの作成
		var label string
		if packet.IsSupported && len(packet.Details) > 0 {
			// サポート対象プロトコルの場合は詳細情報を追加
			details := strings.Join(packet.Details[:minInt(3, len(packet.Details))], "<br>")
			label = fmt.Sprintf("%s: %s", packet.Protocol, details)
		} else {
			// サポート対象外の場合は基本情報のみ
			label = fmt.Sprintf("%s: %s", packet.Protocol, packet.Info)
		}

		// シーケンス図のメッセージ行を追加
		file.WriteString(fmt.Sprintf("    %s->>%s: %s\n", packet.Source, packet.Destination, label))

		// 時間情報を追加（パケットの後に注釈として追加）
		if packet.Time != "" {
			file.WriteString(fmt.Sprintf("    Note over Timeline: %s\n", formatTime(packet.Time)))
		}

		if w.DebugMode {
			fmt.Printf("シーケンス図メッセージ: %s -> %s (%s)\n", packet.Source, packet.Destination, label)
		}
	}

	// Mermaid終了
	file.WriteString("```\n")

	if w.DebugMode {
		fmt.Println("ファイル出力完了:", w.OutputPath)
	}

	return nil
}

// minInt は2つの整数の小さい方を返す補助関数
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// 時刻フォーマットを整形する補助関数
func formatTime(timeStr string) string {
	// tsharkの時刻フォーマットを解析
	t, err := time.Parse("Jan 2, 2006 15:04:05.000000000", timeStr)
	if err != nil {
		// 解析できない場合はそのまま返す
		return timeStr
	}

	// 見やすいフォーマットに変換
	return t.Format("15:04:05.000")
}
