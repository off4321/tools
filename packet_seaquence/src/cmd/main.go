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

	// コマンドライン引数の解析
	filePath := flag.String("file", "", "解析するpcapファイルのパス (必須)")
	outputPath := flag.String("out", "output.md", "出力するマークダウンファイルのパス")
	maxPackets := flag.Int("max", 0, "処理する最大パケット数 (0=すべて)")
	debugMode := flag.Bool("debug", false, "デバッグモードの有効化")
	sourceIP := flag.String("source", "", "送信元IPをフィルタリングする(例: 123.456.789.012)")
	destinationIP := flag.String("destination", "", "送信先IPをフィルタリングする(例: 123.456.789.012)")
	protocolName := flag.String("protocol", "", "プロトコル名をフィルタリングする")
	ipFlag := flag.String("IP", "", "送信元または送信先IPをフィルタリングする(例: 123.456.789.012)")
	startTimeArg := flag.String("startTime", "", "フィルタ適用開始時刻(例: 2023-10-01 12:00:00)")
	endTimeArg := flag.String("endTime", "", "フィルタ適用終了時刻(例: 2023-10-01 12:00:00)")
	version := flag.Bool("version", false, "バージョン情報の表示")

	flag.Parse()

	// バージョン情報の表示
	if *version {
		fmt.Println("Packet Sequence Analyzer v0.0.1")
		fmt.Println("created by Junnosuke Horiuchi 2025-04-06")
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
	fmt.Println("ステップ4: プロトコル詳細情報の解析...")
	for i, packet := range packets {
		if packet.IsSupported {
			if *debugMode {
				fmt.Printf("パケット %d/%d (%s) の詳細情報を取得中...\n", i+1, len(packets), packet.Protocol)
			}
			err := infoGetter.GetDetailedInfo(packet)
			if err != nil {
				fmt.Printf("詳細情報取得エラー (パケット %s): %v\n", packet.Number, err)
				continue
			}
		}
	}
	fmt.Println("詳細情報の解析が完了しました")

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
