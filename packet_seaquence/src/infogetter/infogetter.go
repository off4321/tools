package infogetter

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/packet_sequence/models"
	"github.com/packet_sequence/protocol"
)

// InfoGetter はパケット情報の取得を管理するインターフェース
type InfoGetter interface {
	GetPacketInfo() ([]*models.Packet, error)
	GetDetailedInfo(packet *models.Packet) error
}

// TsharkInfoGetter はtsharkを使用してパケット情報を取得する構造体
type TsharkInfoGetter struct {
	FilePath      string                       // PCAPファイルのパス
	DebugMode     bool                         // デバッグモード
	MaxPackets    int                          // 処理する最大パケット数
	analyzers     map[string]protocol.Analyzer // プロトコル別のアナライザー
	SourceIP      string
	DestinationIP string
	ProtocolName  string
	IPFlag        string
	StartTime     string
	EndTime       string
	packetDetails map[string]string // パケット番号をキーとした詳細情報のマップ
	initialized   bool              // 初期化フラグ
}

// NewTsharkInfoGetter は新しいTsharkInfoGetterを生成する
func NewTsharkInfoGetter(filePath string, debugMode bool, maxPackets int, sourceIP, destinationIP, protocolName, ipFlag, startTime, endTime string) *TsharkInfoGetter {
	getter := &TsharkInfoGetter{
		FilePath:      filePath,
		DebugMode:     debugMode,
		MaxPackets:    maxPackets,
		analyzers:     make(map[string]protocol.Analyzer),
		SourceIP:      sourceIP,
		DestinationIP: destinationIP,
		ProtocolName:  protocolName,
		IPFlag:        ipFlag,
		StartTime:     startTime,
		EndTime:       endTime,
		packetDetails: make(map[string]string),
		initialized:   false,
	}

	// 各プロトコルのアナライザーを初期化
	getter.analyzers["TCP"] = protocol.NewTCPAnalyzer(debugMode)
	getter.analyzers["UDP"] = protocol.NewUDPAnalyzer(debugMode)
	getter.analyzers["ARP"] = protocol.NewARPAnalyzer(debugMode)
	getter.analyzers["ICMP"] = protocol.NewICMPAnalyzer(debugMode)
	getter.analyzers["IPv4"] = protocol.NewIPv4Analyzer(debugMode)
	getter.analyzers["SCTP"] = protocol.NewSCTPAnalyzer(debugMode)
	getter.analyzers["HTTP"] = protocol.NewHTTPAnalyzer(debugMode)
	getter.analyzers["DNS"] = protocol.NewDNSAnalyzer(debugMode)
	getter.analyzers["HTTPS"] = protocol.NewHTTPSAnalyzer(debugMode)
	getter.analyzers["X.25"] = protocol.NewX25Analyzer(debugMode)
	getter.analyzers["NTP"] = protocol.NewNTPAnalyzer(debugMode)

	return getter
}

func readTsharkPath() string {
	if runtime.GOOS == "windows" || runtime.GOOS == "MacOS" {
		// 設定ファイルから読み込み
		execPath, err := os.Executable()
		if err != nil {
			fmt.Printf("実行ファイルのパスを取得できません: %v\n", err)
			return "tshark" // デフォルトのtsharkを返す
		}
		execDir := filepath.Dir(execPath)
		configPath := filepath.Join(execDir, "config", "config.pkseq")
		currentDir, _ := os.Getwd()
		fmt.Printf("カレントディレクトリ: %s\n", currentDir)
		file, err := os.Open(configPath)
		if err != nil {
			fmt.Printf("設定ファイル %s が開けません: %v\n", configPath, err)
			fmt.Println("tsharkが見つかりません。設定ファイルを確認してください。")
			return "tshark" // 設定がなくても"tshark"を返す
		}
		fmt.Println("設定ファイルを読み込み中...")

		defer file.Close()
		reader := bufio.NewReader(file)

		for {
			line, err := reader.ReadString('\n')
			if err == io.EOF || err != nil {
				fmt.Println("設定ファイルの読み込みが完了しました。")
				//break
			}
			line = strings.TrimSpace(line)
			fmt.Println("debug-読み込んだ行:", line)
			if strings.HasPrefix(line, "tsharkDir=") {
				// 引用符を取り除く処理
				dir := strings.TrimPrefix(line, "tsharkDir=")
				dir = strings.Trim(dir, "\"")

				// ファイル名部分を除去してディレクトリパスのみを取得
				tsharkDir := filepath.Dir(dir)
				fmt.Println("debug-tsharkDir:", tsharkDir)

				// PATHに tsharkDir を追加: 現在のPATH + ";" + tsharkDir
				currentPath := os.Getenv("PATH")
				newPath := currentPath + ";" + tsharkDir
				os.Setenv("PATH", newPath)

				//fmt.Println("-----------------------------------------------")
				//fmt.Println("tsharkをPATHに追加しました")
				//fmt.Println("コマンドプロンプトで実行する場合は以下のように入力できます：")
				//fmt.Printf("SET PATH=%%PATH%%;%s\n", tsharkDir)
				//fmt.Println("-----------------------------------------------")

				// tsharkコマンドを返す
				return "tshark"
			}
		}

	} else {
		// Linuxの場合は変更なし
		whichCmd := exec.Command("which", "tshark")
		output, err := whichCmd.Output()
		if err == nil && len(output) > 0 {
			return strings.TrimSpace(string(output))
		}
		return "tshark"
	}
}

// GetPacketInfo はtsharkを使用してパケット情報を取得する
func (g *TsharkInfoGetter) GetPacketInfo() ([]*models.Packet, error) {
	tsharkPath := readTsharkPath()

	// 入力ファイルパスを絶対パスに変換
	absFilePath, err := filepath.Abs(g.FilePath)
	if err != nil {
		fmt.Printf("警告: ファイルパスを絶対パスに変換できません: %v\n", err)
		absFilePath = g.FilePath // 変換失敗時は元のパスを使用
	} else {
		fmt.Printf("入力ファイルの絶対パス: %s\n", absFilePath)
	}

	// ファイルの存在確認
	if _, err := os.Stat(absFilePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("ファイルが見つかりません: %s", absFilePath)
	}

	// tsharkコマンドを構築
	maxOption := ""
	if g.MaxPackets > 0 {
		maxOption = fmt.Sprintf("-c %d", g.MaxPackets)
	}

	filterParts := []string{}
	if g.SourceIP != "" {
		filterParts = append(filterParts, fmt.Sprintf("ip.src == %s", g.SourceIP))
	}
	if g.DestinationIP != "" {
		filterParts = append(filterParts, fmt.Sprintf("ip.dst == %s", g.DestinationIP))
	}
	protocolFilter := strings.ToLower(g.ProtocolName)
	if protocolFilter != "" {
		filterParts = append(filterParts, protocolFilter)
	}
	if g.IPFlag != "" {
		filterParts = append(filterParts, fmt.Sprintf("(ip.src == %s or ip.dst == %s)", g.IPFlag, g.IPFlag))
	}
	if g.StartTime != "" && g.EndTime != "" {
		filterParts = append(filterParts, fmt.Sprintf("frame.time >= \"%s\" and frame.time <= \"%s\"", g.StartTime, g.EndTime))
	}

	var filterExpression string
	if len(filterParts) > 0 {
		filterExpression = strings.Join(filterParts, " and ")
	}

	// 最初に簡易コマンドでファイル形式をチェック
	if runtime.GOOS == "windows" {
		testCmd := exec.Command(tsharkPath, "-r", absFilePath, "-c", "1")
		testOutput, testErr := testCmd.CombinedOutput()
		if testErr != nil {
			if g.DebugMode {
				fmt.Printf("tsharkファイル読み込みテスト失敗: %v\n出力: %s\n", testErr, string(testOutput))
			}
			// より詳細な情報を表示するため、-vオプションを追加したコマンドを試す
			verboseCmd := exec.Command(tsharkPath, "-v")
			verboseOutput, _ := verboseCmd.CombinedOutput()
			fmt.Printf("tsharkバージョン情報:\n%s\n", string(verboseOutput))
		}
	}

	// 最初に概要情報を取得
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" || runtime.GOOS == "MacOS" {
		// WindowsまたはMacOSの場合、tsharkのパスを取得
		args := []string{"-r", absFilePath, "-T", "fields",
			"-e", "frame.number", "-e", "frame.time", "-e", "ip.src", "-e", "ip.dst",
			"-e", "_ws.col.Protocol", "-e", "frame.len", "-e", "_ws.col.Info",
			"-e", "eth.src", "-e", "eth.dst", "-e", "arp.src.proto_ipv4", "-e", "arp.dst.proto_ipv4"}

		if maxOption != "" {
			args = append([]string{"-r", absFilePath, maxOption}, args[2:]...)
		}

		if filterExpression != "" {
			args = append(args, "-Y", filterExpression)
		}

		cmd = exec.Command(tsharkPath, args...)
	} else {
		baseCmd := fmt.Sprintf("tshark -r %s %s -T fields -e frame.number -e frame.time -e ip.src -e ip.dst "+
			"-e _ws.col.Protocol -e frame.len -e _ws.col.Info -e eth.src -e eth.dst "+
			"-e arp.src.proto_ipv4 -e arp.dst.proto_ipv4", g.FilePath, maxOption)
		if filterExpression != "" {
			baseCmd = fmt.Sprintf("tshark -r %s %s -Y '%s' -T fields -e frame.number -e frame.time "+
				"-e ip.src -e ip.dst -e _ws.col.Protocol -e frame.len -e _ws.col.Info "+
				"-e eth.src -e eth.dst -e arp.src.proto_ipv4 -e arp.dst.proto_ipv4",
				g.FilePath, maxOption, filterExpression)
		}
		cmd = exec.Command("sh", "-c", strings.Replace(baseCmd, "tshark", tsharkPath, 1))
	}

	if g.DebugMode {
		fmt.Println("tsharkのパス:", tsharkPath)
		fmt.Println("実行コマンド:", cmd.String())
	}

	// コマンド実行
	output, err := cmd.CombinedOutput() // StdoutとStderrの両方を取得
	if err != nil {
		if g.DebugMode {
			fmt.Printf("tshark出力: %s\n", string(output))
		}
		return nil, fmt.Errorf("tshark実行エラー: %v", err)
	}

	// 出力結果をパースしてパケット情報を構築
	packets := []*models.Packet{}
	scanner := bufio.NewScanner(strings.NewReader(string(output)))

	// パケット番号のリストを保持
	packetNumbers := []string{}

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, "\t")

		// 最低限必要なフィールド数の確認
		if len(fields) < 7 {
			continue // 不正なフォーマット
		}

		number := fields[0]   // Number
		time := fields[1]     // Time
		protocol := fields[4] // Protocol
		length := fields[5]   // Length
		info := fields[6]     // Info

		// IPアドレス (fields[2]とfields[3]) が空の場合の処理
		source := fields[2]
		destination := fields[3]

		// ARPプロトコルの場合は特別な処理
		if protocol == "ARP" {
			// ARPのIPアドレス (fields[8]とfields[9]) を使用
			if len(fields) > 9 {
				if fields[8] != "" {
					source = fields[8] // ARP送信元IPアドレス
				} else if len(fields) > 7 && fields[7] != "" {
					source = fields[7] // MACアドレスを使用
				}

				if fields[9] != "" {
					destination = fields[9] // ARP宛先IPアドレス
				} else if len(fields) > 8 && fields[8] != "" {
					destination = fields[8] // MACアドレスを使用
				}
			}
		}

		// それでも送信元/宛先が空の場合はno_source,no_destinationで代用
		if source == "" && len(fields) > 7 {
			source = "no_source" // 空用のプレースホルダー
		}

		if destination == "" && len(fields) > 8 {
			destination = "no_destination" // 空用のプレースホルダー
		}

		// パケット番号を保存
		packetNumbers = append(packetNumbers, number)

		// パケット情報の生成
		packet := models.NewPacket(
			number,
			time,
			source,
			destination,
			protocol,
			length,
			info,
		)

		packets = append(packets, packet)

		if g.DebugMode {
			fmt.Printf("パケット %s: %s -> %s (%s)\n", packet.Number, packet.Source, packet.Destination, packet.Protocol)
		}
	}

	// 詳細情報を一括で取得し、メモリに格納（初期化がまだの場合）
	if !g.initialized && len(packets) > 0 {
		if g.DebugMode {
			fmt.Printf("詳細情報のプリロードを開始します（パケット数: %d）\n", len(packetNumbers))
		}

		// パケット詳細情報を取得する
		err = g.preloadPacketDetails(packetNumbers, tsharkPath)
		if err != nil {
			if g.DebugMode {
				fmt.Printf("詳細情報のプリロードに失敗: %v\n", err)
			}
		} else {
			if g.DebugMode {
				fmt.Printf("詳細情報のプリロード完了: %d パケットの詳細をキャッシュしました\n", len(g.packetDetails))
			}
			g.initialized = true
		}
	}

	return packets, nil
}

// 詳細情報を一括でプリロードする
func (g *TsharkInfoGetter) preloadPacketDetails(packetNumbers []string, tsharkPath string) error {
	// パケット番号をフィルタとして結合
	if len(packetNumbers) == 0 {
		return nil
	}

	// Linuxでは特に大量のパケットを一度に処理できない可能性があるため、バッチサイズを小さく設定
	batchSize := 20
	if runtime.GOOS == "windows" || runtime.GOOS == "MacOS" {
		batchSize = 100
	}

	totalBatches := (len(packetNumbers) + batchSize - 1) / batchSize

	if g.DebugMode {
		fmt.Printf("詳細情報のプリロード: %d パケットを %d バッチに分割します\n", len(packetNumbers), totalBatches)
	}

	for i := 0; i < totalBatches; i++ {
		start := i * batchSize
		end := (i + 1) * batchSize
		if end > len(packetNumbers) {
			end = len(packetNumbers)
		}

		batchNumbers := packetNumbers[start:end]
		frameFilters := make([]string, len(batchNumbers))
		for j, num := range batchNumbers {
			frameFilters[j] = fmt.Sprintf("frame.number==%s", num)
		}

		filterExpression := strings.Join(frameFilters, " or ")

		var cmd *exec.Cmd
		if runtime.GOOS == "windows" || runtime.GOOS == "MacOS" {
			cmd = exec.Command(tsharkPath, "-r", g.FilePath, "-Y", filterExpression, "-V")
		} else {
			// Linuxの場合、コマンドを慎重に構築
			cmdStr := fmt.Sprintf("%s -r %s -Y \"%s\" -V", tsharkPath, g.FilePath, filterExpression)
			cmd = exec.Command("sh", "-c", cmdStr)
		}

		if g.DebugMode {
			fmt.Printf("プリロードバッチ %d/%d: パケット %d-%d のプリロードを実行\n",
				i+1, totalBatches, start+1, end)
			fmt.Printf("コマンド: %s\n", cmd.String())
		}

		output, err := cmd.Output()
		if err != nil {
			return fmt.Errorf("プリロード実行エラー (バッチ %d/%d): %v", i+1, totalBatches, err)
		}

		// 出力を各パケットごとに分割して保存
		countBefore := len(g.packetDetails)
		g.parsePacketDetails(string(output))
		countAfter := len(g.packetDetails)

		if g.DebugMode {
			fmt.Printf("バッチ %d/%d 完了: キャッシュに %d パケットを追加 (合計: %d)\n",
				i+1, totalBatches, countAfter-countBefore, countAfter)
		}
		fmt.Printf("詳細情報のプリロード進捗: %d/%d バッチ完了\n", i+1, totalBatches)
	}

	return nil
}

// parsePacketDetails は出力を各パケットごとに分割して保存するヘルパーメソッド
func (g *TsharkInfoGetter) parsePacketDetails(output string) {
	// "Frame X: " をセパレータとしてパケットを分割
	frameRegex := regexp.MustCompile(`Frame (\d+):`)
	matches := frameRegex.FindAllStringSubmatchIndex(output, -1)

	if len(matches) == 0 {
		if g.DebugMode {
			fmt.Println("警告: 詳細出力からパケットフレーム情報を検出できませんでした")
			if len(output) > 100 {
				fmt.Printf("出力サンプル: %s...\n", output[:100])
			} else if len(output) > 0 {
				fmt.Printf("出力サンプル: %s\n", output)
			}
		}
		return
	}

	for i := 0; i < len(matches); i++ {
		start := matches[i][0]
		end := len(output)
		if i < len(matches)-1 {
			end = matches[i+1][0]
		}

		frameSection := output[start:end]
		frameMatch := frameRegex.FindStringSubmatch(frameSection)
		if len(frameMatch) > 1 {
			frameNumber := frameMatch[1]
			g.packetDetails[frameNumber] = frameSection
			if g.DebugMode && i == 0 {
				fmt.Printf("パケット %s の詳細情報をキャッシュに保存しました\n", frameNumber)
			}
		}
	}
}

// GetDetailedInfo はプロトコル別の詳細情報を取得する
func (g *TsharkInfoGetter) GetDetailedInfo(packet *models.Packet) error {
	var detailOutput string

	if g.DebugMode {
		fmt.Printf("パケット %s/%s (%s) の詳細情報を取得中...\n",
			packet.Number, packet.Protocol, packet.Protocol)
	}

	// メモリにキャッシュされた詳細情報を確認
	if details, exists := g.packetDetails[packet.Number]; exists {
		if g.DebugMode {
			fmt.Printf("パケット %s の詳細情報はキャッシュから取得しました\n", packet.Number)
		}
		detailOutput = details
	} else {
		// キャッシュにない場合のみ個別にtsharkを実行
		if g.DebugMode {
			fmt.Printf("警告: パケット %s の詳細情報がキャッシュにありません。個別に取得します。\n", packet.Number)
		}

		filter := fmt.Sprintf("frame.number==%s", packet.Number)
		tsharkPath := readTsharkPath()

		var cmd *exec.Cmd
		if runtime.GOOS == "windows" || runtime.GOOS == "MacOS" {
			cmd = exec.Command(tsharkPath, "-r", g.FilePath, "-Y", filter, "-V")
		} else {
			cmd = exec.Command("sh", "-c", fmt.Sprintf("%s -r %s -Y '%s' -V", tsharkPath, g.FilePath, filter))
		}

		if g.DebugMode {
			fmt.Println("詳細情報個別取得コマンド:", cmd.String())
		}

		output, err := cmd.Output()
		if err != nil {
			return fmt.Errorf("詳細情報取得エラー: %v", err)
		}
		detailOutput = string(output)

		// 今後の使用のためにキャッシュに保存
		g.packetDetails[packet.Number] = detailOutput
	}

	// 出力結果から必要な情報を抽出
	scanner := bufio.NewScanner(strings.NewReader(detailOutput))
	protocolSection := false

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		// プロトコルセクションの検出
		if strings.Contains(line, packet.Protocol) {
			protocolSection = true
			continue
		}

		if protocolSection && line != "" && !strings.HasPrefix(line, "Frame") && !strings.HasPrefix(line, "Ethernet") {
			// 詳細情報を追加（余分な空白を削除）
			re := regexp.MustCompile(`\s{2,}`)
			cleanLine := re.ReplaceAllString(line, " ")
			packet.AddDetail(cleanLine)
		}
	}

	// ARPプロトコルの場合、送信元と宛先のアドレスを詳細情報から補完
	if packet.Protocol == "ARP" && (packet.Source == "" || packet.Destination == "") {
		srcIPRe := regexp.MustCompile(`Sender IP address: ([0-9.]+)`)
		dstIPRe := regexp.MustCompile(`Target IP address: ([0-9.]+)`)
		srcMACRe := regexp.MustCompile(`Sender MAC address: ([0-9A-Fa-f:]+)`)
		dstMACRe := regexp.MustCompile(`Target MAC address: ([0-9A-Fa-f:]+)`)

		// 出力結果を再度スキャンして送信元と宛先のアドレスを探す
		scanner = bufio.NewScanner(strings.NewReader(detailOutput))
		for scanner.Scan() {
			line := scanner.Text()
			line = strings.TrimSpace(line)

			// 送信元IPアドレス
			if matches := srcIPRe.FindStringSubmatch(line); len(matches) > 1 {
				if packet.Source == "" {
					packet.Source = matches[1]
				}
			}

			// 宛先IPアドレス
			if matches := dstIPRe.FindStringSubmatch(line); len(matches) > 1 {
				if packet.Destination == "" {
					packet.Destination = matches[1]
				}
			}

			// 送信元MACアドレス (IPが見つからない場合のフォールバック)
			if packet.Source == "" {
				if matches := srcMACRe.FindStringSubmatch(line); len(matches) > 1 {
					packet.Source = matches[1]
				}
			}

			// 宛先MACアドレス (IPが見つからない場合のフォールバック)
			if packet.Destination == "" {
				if matches := dstMACRe.FindStringSubmatch(line); len(matches) > 1 {
					packet.Destination = matches[1]
				}
			}
		}
	}

	// 対応するプロトコルアナライザーを取得
	analyzer, exists := g.analyzers[packet.Protocol]
	if !exists {
		if g.DebugMode {
			fmt.Printf("プロトコル %s のアナライザーがありません\n", packet.Protocol)
		}
		return nil
	}

	// アナライザーを使用して詳細情報を解析
	details, err := analyzer.Analyze(packet)
	if err != nil {
		if g.DebugMode {
			fmt.Printf("プロトコル %s の解析エラー: %v\n", packet.Protocol, err)
		}
		return nil
	}

	// パケットの詳細情報をクリアして解析結果で置き換え
	packet.Details = details

	if g.DebugMode {
		fmt.Printf("パケット %s (%s) の詳細情報:\n", packet.Number, packet.Protocol)
		for _, detail := range details {
			fmt.Printf("  - %s\n", detail)
		}
	}

	return nil
}
