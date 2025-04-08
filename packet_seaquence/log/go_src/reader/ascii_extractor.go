package reader

import (
	"strings"
)

// extractAsciiText extracts printable ASCII text from raw binary data
func extractAsciiText(data []byte) string {
	if len(data) == 0 {
		return ""
	}
	
	var builder strings.Builder
	inTextSequence := false
	minSequenceLen := 3 // 最低3文字以上が連続しないとテキストとみなさない
	currentSequence := ""
	
	for _, b := range data {
		// 印字可能なASCII文字か改行/タブかどうか判定
		if (b >= 32 && b <= 126) || b == 10 || b == 13 || b == 9 {
			inTextSequence = true
			currentSequence += string(b)
		} else {
			// ASCII文字でない場合
			if inTextSequence {
				// シーケンスが最低長を超えていれば追加
				if len(currentSequence) >= minSequenceLen {
					builder.WriteString(currentSequence)
					builder.WriteString(" ")
				}
				inTextSequence = false
				currentSequence = ""
			}
		}
	}
	
	// 最後のシーケンスを処理
	if inTextSequence && len(currentSequence) >= minSequenceLen {
		builder.WriteString(currentSequence)
	}
	
	result := builder.String()
	result = strings.TrimSpace(result)
	
	// 制御文字の置換
	result = strings.ReplaceAll(result, string([]byte{10}), "\\n")
	result = strings.ReplaceAll(result, string([]byte{13}), "\\r")
	result = strings.ReplaceAll(result, string([]byte{9}), "\\t")
	
	return result
}