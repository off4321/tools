// Package main provides the command-line interface for packet sequence analysis
package main

import (
	"flag"
	"fmt"
	"os"
	
	"app/go_src/analyzer"
)

func main() {
	// Define command-line flags
	pcapFile := flag.String("file", "", "PCAPファイルのパス (必須)")
	outputFile := flag.String("output", "output.md", "出力ファイルのパス (デフォルト: output.md)")
	maxEntries := flag.Int("max", 50, "1つのファイルに出力する最大パケット数 (デフォルト: 50)")
	verbose := flag.Bool("verbose", false, "詳細なデバッグ情報を表示する")
	
	// Parse command-line flags
	flag.Parse()
	
	// Check if required flags are provided
	if *pcapFile == "" {
		fmt.Println("エラー: PCAPファイルのパスを指定してください (-file)")
		flag.Usage()
		os.Exit(1)
	}
	
	// Create options map
	options := map[string]interface{}{
		"verbose": *verbose,
	}
	
	// Create the analyzer
	analyzer, err := analyzer.NewPacketSequenceAnalyzer(*pcapFile, *outputFile, *maxEntries, options)
	if err != nil {
		fmt.Printf("エラー: %v\n", err)
		os.Exit(1)
	}
	
	// Run the analysis
	fmt.Printf("PCAPファイル %s を解析中...\n", *pcapFile)
	if err := analyzer.Analyze(); err != nil {
		fmt.Printf("エラー: %v\n", err)
		os.Exit(1)
	}
}