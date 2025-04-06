package protocol

import (
	"github.com/packet_sequence/models"
)

// Analyzer はプロトコル解析機能を提供するインターフェース
type Analyzer interface {
	// Analyze はパケットから詳細情報を抽出して返す
	Analyze(packet *models.Packet) ([]string, error)
}