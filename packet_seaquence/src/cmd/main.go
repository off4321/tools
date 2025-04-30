package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/packet_sequence/checker"
	"github.com/packet_sequence/infogetter"
	"github.com/packet_sequence/models"
	"github.com/packet_sequence/reader"
	"github.com/packet_sequence/writer"
)

func main() {
	// 処理開始時間を記録
	startTime := time.Now()

	// カスタムUsage設定
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "使い方: %s [オプション]\n\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "オプション:")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "\n例:")
		fmt.Fprintln(os.Stderr, "  パケット解析: packet_seaquence -file=capture.pcap -out=result.md")
		fmt.Fprintln(os.Stderr, "  特定IPのみ解析: packet_seaquence -file=capture.pcap -IP=192.168.1.1")
		fmt.Fprintln(os.Stderr, "  詳細モード: packet_seaquence -file=capture.pcap -debug=true")
		fmt.Fprintln(os.Stderr, "  全詳細情報表示: packet_seaquence -file=capture.pcap -info=all")
	}

	// コマンドライン引数の解析
	filePath := flag.String("file", "", "解析するpcapファイルのパス (必須)!!絶対パスを使ってください!!")
	outputPath := flag.String("out", "output.md", "出力するマークダウンファイルのパス")
	maxPackets := flag.Int("max", 0, "処理する最大パケット数 (0=すべて)")
	debugMode := flag.Bool("debug", false, "デバッグモードの有効化")
	sourceIP := flag.String("source", "", "送信元IPをフィルタリングする(例: 123.456.789.012)")
	destinationIP := flag.String("destination", "", "送信先IPをフィルタリングする(例: 123.456.789.012)")
	protocolName := flag.String("protocol", "", "プロトコル名をフィルタリングする")
	ipFlag := flag.String("IP", "", "送信元または送信先IPをフィルタリングする(例: 123.456.789.012)")
	startTimeArg := flag.String("startTime", "", "フィルタ適用開始時刻(例: 2023-10-01 12:00:00)")
	endTimeArg := flag.String("endTime", "", "フィルタ適用終了時刻(例: 2023-10-01 12:00:00)")
	infoArg := flag.String("info", "", "詳細情報の表示オプション(例: all=全詳細情報を表示)")
	version := flag.Bool("version", false, "バージョン情報の表示")
	help := flag.Bool("help", false, "ヘルプ情報の表示")

	flag.Parse()

	// ヘルプの表示
	if *help {
		flag.Usage()
		return
	}

	// バージョン情報の表示
	if *version {
		fmt.Println("Packet Sequence Analyzer v0.0.2")
		fmt.Println("created by Junnosuke Horiuchi 2025-05-01")
		fmt.Println("Github URL:https://github.com/off4321/tools/tree/main/packet_seaquence")
		fmt.Println("License: MIT")
		fmt.Println("Copyright (c) 2023 Junnosuke Horiuchi")
		fmt.Println("This program requires tshark to run.")
		return
	}

	// 必須パラメータのチェック
	if *filePath == "" {
		fmt.Println("エラー: -file パラメータが必要です")
		flag.Usage()
		os.Exit(1)
	}

	// デバッグ情報の表示
	if *debugMode {
		fmt.Println("デバッグモード: 有効")
		fmt.Println("入力ファイル:", *filePath)
		fmt.Println("出力ファイル:", *outputPath)
		fmt.Println("最大パケット数:", *maxPackets)
	}

	// 1. PCAPファイルの読み込み
	fmt.Println("ステップ1: PCAPファイルの読み込み...")
	pcapReader := reader.NewPCAPReader(*filePath)
	err := pcapReader.ReadFile()
	if err != nil {
		fmt.Printf("ファイル読み込みエラー: %v\n", err)
		os.Exit(1)
	}

	// 2. tsharkを使ってpcapファイルの情報取得
	fmt.Println("ステップ2: パケット情報の取得...")

	// infoArgが "all" のときだけInfoAllをtrueにする
	infoAll := *infoArg == "all"
	if infoAll && *debugMode {
		fmt.Println("全詳細情報表示モード: 有効")
	}

	infoGetter := infogetter.NewTsharkInfoGetter(
		*filePath,
		*debugMode,
		*maxPackets,
		*sourceIP,
		*destinationIP,
		*protocolName,
		*ipFlag,
		*startTimeArg,
		*endTimeArg,
		infoAll,
	)
	packets, err := infoGetter.GetPacketInfo()
	if err != nil {
		fmt.Printf("パケット情報取得エラー: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("パケット数: %d\n", len(packets))

	// 3. プロトコルのサポート対象チェック
	fmt.Println("ステップ3: プロトコルのサポート状況をチェック...")
	protocolChecker := checker.NewTsharkProtocolChecker(models.SupportedProtocols, *debugMode)
	protocolChecker.CheckPackets(packets)

	// サポート対象プロトコルの数をカウント
	supportedCount := 0
	for _, packet := range packets {
		if packet.IsSupported {
			supportedCount++
		}
	}
	fmt.Printf("サポート対象プロトコル: %d/%d\n", supportedCount, len(packets))

	// 4. サポート対象プロトコルの詳細情報取得
	err = getPacketDetails(infoGetter, packets)
	if err != nil {
		fmt.Printf("詳細情報取得エラー: %v\n", err)
		os.Exit(1)
	}

	// 5. Markdown形式での出力
	fmt.Println("ステップ5: シーケンス図の生成...")
	mermaidWriter := writer.NewMermaidWriter(*outputPath, *debugMode)
	err = mermaidWriter.Write(packets)
	if err != nil {
		fmt.Printf("ファイル出力エラー: %v\n", err)
		os.Exit(1)
	}

	// 処理時間を計算
	elapsedTime := time.Since(startTime)
	fmt.Printf("\n処理が完了しました。出力ファイル: %s (所要時間: %v)\n", *outputPath, elapsedTime)
}

// getPacketDetails はパケットの詳細情報を取得する
func getPacketDetails(infoGetter infogetter.InfoGetter, packets []*models.Packet) error {
	fmt.Printf("ステップ4: プロトコル詳細情報の解析...\n")

	// infoGetterをTsharkInfoGetterにキャストしてデバッグモードを取得
	tsharkGetter, ok := infoGetter.(*infogetter.TsharkInfoGetter)
	isDebug := false
	if ok {
		isDebug = tsharkGetter.DebugMode
	}

	// この時点で既に詳細情報はプリロードされているはず
	for i, packet := range packets {
		if isDebug {
			fmt.Printf("パケット %d/%d (%s) の詳細情報を処理中...\n", i+1, len(packets), packet.Protocol)
		} else if (i+1)%1000 == 0 || i+1 == len(packets) {
			// デバッグモードでなければ1000パケットごと、または最後のパケットのみ進捗を表示
			fmt.Printf("パケット処理進捗: %d/%d\n", i+1, len(packets))
		}

		err := infoGetter.GetDetailedInfo(packet)
		if err != nil {
			// 詳細情報の取得に失敗した場合はエラーメッセージを表示
			// ただし、処理は続行する
			fmt.Printf("詳細情報取得エラー: %v\n", err)
		}
	}

	fmt.Println("詳細情報の解析が完了しました")

	return nil
}
