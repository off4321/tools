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
	// デバッグのために現在の作業ディレクトリを表示
	//currentDir, _ := os.Getwd()
	//mt.Println("現在の作業ディレクトリ:", currentDir)

	if runtime.GOOS == "windows" || runtime.GOOS == "MacOS" {
		// 設定ファイルから読み込み
		path := "config/config.pkseq"
		//fmt.Println("設定ファイル検索:", path)
		file, err := os.Open(path)
		if err == nil {
			//fmt.Println("設定ファイルを見つけました:", path)
			defer file.Close()
			reader := bufio.NewReader(file)
			for {
				line, err := reader.ReadString('\n')
				if err == io.EOF || err != nil {
					break
				}
				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, "tsharkDir=") {
					// 引用符を取り除く処理
					path := strings.TrimPrefix(line, "tsharkDir=")
					path = strings.Trim(path, "\"")
					//fmt.Println("tsharkパスを設定:", path)
					return path
				}
			}
		}

		// 環境変数PATHからtsharkを探す
		pathEnv := os.Getenv("PATH")
		pathDirs := strings.Split(pathEnv, ";") // Windowsの場合

		for _, dir := range pathDirs {
			tsharkPath := filepath.Join(dir, "tshark.exe")
			if _, err := os.Stat(tsharkPath); err == nil {
				//fmt.Println("PATHから見つかったtshark:", tsharkPath)
				return tsharkPath
			}
		}

		// Wiresharkの標準インストールディレクトリを確認
		commonPaths := []string{
			"C:\\Program Files\\Wireshark\\tshark.exe",
			"C:\\Program Files (x86)\\Wireshark\\tshark.exe",
		}

		for _, path := range commonPaths {
			if _, err := os.Stat(path); err == nil {
				//fmt.Println("標準パスから見つかったtshark:", path)
				return path
			}
		}

		fmt.Println("tsharkが見つかりません。デフォルトのパスを使用します。")
		return "tshark"
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

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" || runtime.GOOS == "MacOS" {
		// WindowsまたはMacOSの場合、tsharkのパスを取得
		if filterExpression != "" {
			cmd = exec.Command(tsharkPath, "-r", g.FilePath, maxOption, "-Y", filterExpression, "-T", "fields",
				"-e", "frame.number", "-e", "frame.time", "-e", "ip.src", "-e", "ip.dst",
				"-e", "_ws.col.Protocol", "-e", "frame.len", "-e", "_ws.col.Info",
				"-e", "eth.src", "-e", "eth.dst", "-e", "arp.src.proto_ipv4", "-e", "arp.dst.proto_ipv4")
		} else {
			cmd = exec.Command(tsharkPath, "-r", g.FilePath, maxOption, "-T", "fields",
				"-e", "frame.number", "-e", "frame.time", "-e", "ip.src", "-e", "ip.dst",
				"-e", "_ws.col.Protocol", "-e", "frame.len", "-e", "_ws.col.Info",
				"-e", "eth.src", "-e", "eth.dst", "-e", "arp.src.proto_ipv4", "-e", "arp.dst.proto_ipv4")
		}
	} else {
		// Linux等で従来の"sh -c"呼び出し
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
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("tshark実行エラー: %v", err)
	}

	// 出力結果をパースしてパケット情報を構築
	packets := []*models.Packet{}
	scanner := bufio.NewScanner(strings.NewReader(string(output)))

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

	return packets, nil
}

// GetDetailedInfo はプロトコル別の詳細情報を取得する
func (g *TsharkInfoGetter) GetDetailedInfo(packet *models.Packet) error {
	// プロトコル名に応じたtsharkフィルタを構築
	filter := fmt.Sprintf("frame.number==%s", packet.Number)
	tsharkPath := readTsharkPath()

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" || runtime.GOOS == "MacOS" {
		// WindowsまたはMacOSの場合、直接コマンドを実行
		cmd = exec.Command(tsharkPath, "-r", g.FilePath, "-Y", filter, "-V")
	} else {
		// Linux等では従来の"sh -c"呼び出し
		cmd = exec.Command("sh", "-c", fmt.Sprintf("%s -r %s -Y '%s' -V", tsharkPath, g.FilePath, filter))
	}

	if g.DebugMode {
		fmt.Println("詳細情報取得コマンド:", cmd.String())
	}

	// コマンド実行
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("詳細情報取得エラー: %v", err)
	}

	// 出力結果から必要な情報を抽出
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
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
		scanner = bufio.NewScanner(strings.NewReader(string(output)))
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
