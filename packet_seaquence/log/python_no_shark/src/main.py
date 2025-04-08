#!/usr/bin/env python3
"""
packet_sequenceツールのメインエントリーポイント
"""
import sys
import os
import argparse
from src.analyzer import PacketSequenceAnalyzer


def parse_arguments():
    """
    コマンドライン引数を解析する
    
    Returns:
        argparse.Namespace: 解析された引数
    """
    parser = argparse.ArgumentParser(description='PCAPファイルからシーケンス図を生成するツール')
    
    parser.add_argument('pcap_file', nargs='?', help='PCAPファイルのパス')
    parser.add_argument('-i', '--input', help='入力ファイルのパス（pcap_fileの代わりに指定可能）')
    parser.add_argument('-o', '--output', default='output.mmd', help='出力ファイルのパス（デフォルト: output.mmd）')
    parser.add_argument('-f', '--format', choices=['mermaid', 'text'], default='mermaid', help='出力フォーマット（デフォルト: mermaid）')
    parser.add_argument('-m', '--max-entries', type=int, default=50, help='抽出するエントリの最大数（デフォルト: 50）')
    parser.add_argument('-v', '--verbose', action='store_true', help='詳細な出力を有効にする')
    parser.add_argument('-version', action='store_true', help='ツールのバージョンを表示する')
    
    args = parser.parse_args()
    
    # バージョン表示の場合
    if args.version:
        print("packet_sequence v0.1.0")
        sys.exit(0)
    
    # 入力ファイルのパスを決定
    if args.pcap_file:
        args.input = args.pcap_file
    elif not args.input:
        parser.error("PCAPファイルのパスを指定してください")
    
    # ファイルの存在確認
    if not os.path.exists(args.input):
        parser.error(f"指定されたPCAPファイルが存在しません: {args.input}")
    
    return args


def main():
    """
    メインエントリーポイント
    """
    try:
        # コマンドライン引数を解析
        args = parse_arguments()
        
        # オプション設定
        options = {
            'verbose': args.verbose,
            'format': args.format
        }
        
        # 解析を実行
        analyzer = PacketSequenceAnalyzer(
            pcap_file=args.input,
            output_file=args.output,
            max_entries=args.max_entries,
            options=options
        )
        
        success = analyzer.analyze()
        
        # 結果に基づいて終了コードを設定
        sys.exit(0 if success else 1)
        
    except KeyboardInterrupt:
        print("\n処理が中断されました")
        sys.exit(130)
    except Exception as e:
        print(f"エラーが発生しました: {str(e)}")
        import traceback
        traceback.print_exc()
        sys.exit(1)


if __name__ == "__main__":
    main()
