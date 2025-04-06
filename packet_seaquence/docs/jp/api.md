# API ドキュメント
このドキュメントは、ツールの内部動作を理解し、拡張または変更したい開発者向けに、`関数の説明`、`パラメータ`、`戻り値`を含む`パケットシーケンスツールのAPI参照`を提供します。

## パッケージ構成

### models パッケージ
パケット情報のデータモデルを定義します。

```go
// Packet はパケット情報を格納する基本構造体
type Packet struct {
	Number      string   // パケット番号
	Time        string   // パケットの捕捉時間
	Source      string   // 送信元アドレス
	Destination string   // 宛先アドレス
	Protocol    string   // プロトコル名
	Length      string   // パケット長
	Info        string   // 基本情報
	Details     []string // 詳細情報
	IsSupported bool     // サポート対象プロトコルかどうか
}

// NewPacket は新しいPacketインスタンスを生成する
func NewPacket(number, time, source, destination, protocol, length, info string) *Packet

// AddDetail はパケットに詳細情報を追加する
func (p *Packet) AddDetail(detail string)
```

### reader パッケージ
PCAPファイルの読み込みを管理します。

```go
// Reader はパケットファイルの読み込みを管理するインターフェース
type Reader interface {
	ReadFile() error
	GetFilePath() string
	Exists() bool
}

// NewPCAPReader は新しいPCAPReaderを生成する
func NewPCAPReader(filePath string) *PCAPReader
```

### infogetter パッケージ
パケット情報の取得を管理します。

```go
// InfoGetter はパケット情報の取得を管理するインターフェース
type InfoGetter interface {
	GetPacketInfo() ([]*models.Packet, error)
	GetDetailedInfo(packet *models.Packet) error
}

// NewTsharkInfoGetter は新しいTsharkInfoGetterを生成する
func NewTsharkInfoGetter(filePath string, debugMode bool, maxPackets int, sourceIP, destinationIP, protocolName, ipFlag, startTime, endTime string) *TsharkInfoGetter
```

### protocol パッケージ
各プロトコルの解析機能を提供します。

```go
// Analyzer はプロトコル解析機能を提供するインターフェース
type Analyzer interface {
	Analyze(packet *models.Packet) ([]string, error)
}

// 各プロトコルアナライザーの生成関数
func NewTCPAnalyzer(debugMode bool) *TCPAnalyzer
func NewUDPAnalyzer(debugMode bool) *UDPAnalyzer
func NewARPAnalyzer(debugMode bool) *ARPAnalyzer
func NewICMPAnalyzer(debugMode bool) *ICMPAnalyzer
func NewDNSAnalyzer(debugMode bool) *DNSAnalyzer
func NewHTTPAnalyzer(debugMode bool) *HTTPAnalyzer
func NewHTTPSAnalyzer(debugMode bool) *HTTPSAnalyzer
func NewSCTPAnalyzer(debugMode bool) *SCTPAnalyzer
func NewIPv4Analyzer(debugMode bool) *IPv4Analyzer
func NewX25Analyzer(debugMode bool) *X25Analyzer
func NewNTPAnalyzer(debugMode bool) *NTPAnalyzer
```

### writer パッケージ
パケット情報を出力する機能を提供します。

```go
// Writer はパケット情報を出力するインターフェース
type Writer interface {
	Write(packets []*models.Packet) error
}

// NewMermaidWriter は新しいMermaidWriterを生成する
func NewMermaidWriter(outputPath string, debugMode bool) *MermaidWriter
```

### checker パッケージ
プロトコルのサポート状況を確認します。

```go
// ProtocolChecker はプロトコルがサポート対象かどうかを判定するインターフェース
type ProtocolChecker interface {
	IsSupported(protocol string) bool
	CheckPackets(packets []*models.Packet)
}

// NewTsharkProtocolChecker は新しいTsharkProtocolCheckerを生成する
func NewTsharkProtocolChecker(supportedProtocols []string, debugMode bool) *TsharkProtocolChecker
```

## メインプログラムのフロー
1. コマンドライン引数の解析
2. PCAPファイルの読み込み
3. パケット情報の取得
4. プロトコルのサポート状況確認
5. 詳細情報の取得と解析
6. シーケンス図の生成と出力
