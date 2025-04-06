package reader

import (
	"errors"
	"os"
)

// Reader はパケットファイルの読み込みを管理するインターフェース
type Reader interface {
	ReadFile() error
	GetFilePath() string
	Exists() bool
}

// PCAPReader はPCAPファイルの読み込みを実装する構造体
type PCAPReader struct {
	FilePath string // PCAPファイルのパス
}

// NewPCAPReader は新しいPCAPReaderを生成する
func NewPCAPReader(filePath string) *PCAPReader {
	return &PCAPReader{
		FilePath: filePath,
	}
}

// ReadFile はファイルが読み込み可能かを確認します
func (r *PCAPReader) ReadFile() error {
	if !r.Exists() {
		return errors.New("ファイルが存在しません: " + r.FilePath)
	}
	return nil
}

// GetFilePath はファイルパスを返します
func (r *PCAPReader) GetFilePath() string {
	return r.FilePath
}

// Exists はファイルが存在するかを確認します
func (r *PCAPReader) Exists() bool {
	_, err := os.Stat(r.FilePath)
	return err == nil
}