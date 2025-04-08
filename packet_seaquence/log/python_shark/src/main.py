#!/usr/bin/env python3
"""
パケットシーケンス分析ツールのメインエントリポイント
"""
import os
import sys
import stat

# 重要：他のインポートより前に配置
def setup_bundled_tshark():
    """同梱されたtsharkを設定"""
    # 実行環境がコンパイル済みかどうか判定
    is_frozen = getattr(sys, 'frozen', False)
    
    if is_frozen:
        # コンパイル済み環境の場合
        base_dir = os.path.dirname(sys.executable)
        temp_dir = os.path.join(base_dir, 'temp')
        os.makedirs(temp_dir, exist_ok=True)
        
        bundled_tshark = os.path.join(base_dir, 'tshark')
        if os.path.exists(bundled_tshark):
            # 実行権限を付与
            try:
                os.chmod(bundled_tshark, stat.S_IRWXU | stat.S_IRGRP | stat.S_IXGRP)
                
                # 環境変数を設定
                os.environ["PYSHARK_TSHARK_PATH"] = bundled_tshark
                
                # wireshark_libsが存在する場合、LD_LIBRARY_PATHを設定
                wireshark_libs = os.path.join(base_dir, 'wireshark_libs')
                if os.path.exists(wireshark_libs):
                    if 'LD_LIBRARY_PATH' in os.environ:
                        os.environ['LD_LIBRARY_PATH'] = f"{wireshark_libs}:{os.environ['LD_LIBRARY_PATH']}"
                    else:
                        os.environ['LD_LIBRARY_PATH'] = wireshark_libs
                
                print(f"同梱されたtsharkを使用: {bundled_tshark}")
                return True
            except Exception as e:
                print(f"同梱tsharkの設定中にエラー: {e}")
    
    # システムのtsharkを探す (フォールバック)
    for path in ['/usr/bin/tshark', '/usr/local/bin/tshark']:
        if os.path.exists(path) and os.access(path, os.X_OK):
            os.environ["PYSHARK_TSHARK_PATH"] = path
            print(f"システムのtsharkを使用: {path}")
            return True
    
    print("警告: tsharkが見つかりませんでした")
    return False

# tsharkの設定を実行 (他のインポートより前に)
setup_bundled_tshark()

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
