#!/usr/bin/env python3
"""
PacketSequenceAnalyzerクラスを提供するモジュール
"""
import os
from .reader import ReadPcap
from .extractor import ExtractData
from .writer.mermaid import WriteMermaid
from .writer.file import WriteFile
from .models import PacketSequenceData


class PacketSequenceAnalyzer:
    """
    パケットシーケンス解析の全体の流れを管理するクラス
    
    設計書に記載されている通り、以下の流れで処理を行う:
    1. PCAPファイルの読み込み
    2. データの抽出
    3. Mermaid記法によるシーケンス図の生成
    4. ファイル出力
    """
    
    def __init__(self, pcap_file, output_file="output.mmd", max_entries=50, options=None):
        """
        PacketSequenceAnalyzerクラスの初期化
        
        Args:
            pcap_file (str): PCAPファイルのパス
            output_file (str): 出力ファイルのパス (デフォルト: output.mmd)
            max_entries (int): 抽出するエントリの最大数 (デフォルト: 50)
            options (dict): オプション設定
        """
        self.pcap_file = pcap_file
        self.output_file = output_file
        self.max_entries = max_entries
        self.options = options if options else {}
        
        # ファイルの存在確認
        if not os.path.exists(pcap_file):
            raise FileNotFoundError(f"PCAPファイルが見つかりません: {pcap_file}")
    
    def analyze(self):
        """
        パケットシーケンス解析を実行する
        
        処理の流れ:
        1. PCAPファイルを読み込む
        2. パケットからデータを抽出する
        3. Mermaid記法でシーケンス図を生成する
        4. 最大パケット数に基づいてファイルに分割出力する
        
        Returns:
            bool: 解析が成功したかどうか
        """
        try:
            # PCAPファイルの読み込み
            reader = ReadPcap(self.pcap_file)
            packets = reader.read()
            
            # データの抽出
            extractor = ExtractData(packets, len(packets))  # 全パケットを抽出
            packet_data = extractor.extract()
            
            # PacketSequenceDataオブジェクトの詳細なデバッグ出力
            if 'verbose' in self.options and self.options['verbose']:
                print("\n===== PacketSequenceData DEBUG =====")
                print(f"Type: {type(packet_data)}")
                
                # パケット一覧を取得して内容をデバッグ表示
                packets_list = packet_data.get_packets()
                print(f"Total packets: {len(packets_list)}")
                
                # 最初の数パケットの詳細を表示
                for i, packet in enumerate(packets_list[:3]):  # 最初の3パケットだけ表示
                    print(f"\nPacket {i+1}:")
                    print(f"  Protocol: {packet.get('protocol', 'UNKNOWN')}")
                    print(f"  Src: {packet.get('src', 'Unknown')}")
                    print(f"  Dst: {packet.get('dst', 'Unknown')}")
                    print(f"  Time: {packet.get('time', 'Unknown')}")
                    
                    # infoの内容を表示
                    info = packet.get('info', {})
                    print(f"  Info keys: {list(info.keys())}")
                    
                    # 各プロトコル情報の詳細を表示
                    for key in info:
                        print(f"  {key}: {info[key]}")
                
                print("=================================\n")
                
                # WriteMermaidのbuild_messageメソッドをデバッグ
                print("\n===== WriteMermaid DEBUG =====")
                mermaid_writer = WriteMermaid(packet_data.copy_first_n(3))
                # 最初の数パケットでbuild_messageをテスト
                packets_list = packet_data.get_packets()
                for i, packet in enumerate(packets_list[:3]):
                    message = mermaid_writer._build_message(packet)
                    print(f"Packet {i+1} message: '{message}'")
                print("=================================\n")
            
            # パケットリストを取得
            all_packets = packet_data.get_packets()
            total_packets = len(all_packets)
            
            # パケット数をmax_entriesで分割
            chunk_size = self.max_entries
            chunks = []
            
            for i in range(0, total_packets, chunk_size):
                chunk = all_packets[i:i+chunk_size]
                # PacketSequenceDataオブジェクトを正しくインポートされた状態で作成
                chunk_data = PacketSequenceData()
                for packet in chunk:
                    chunk_data.add_packet(packet)
                chunks.append(chunk_data)
            
            # 各チャンクに対してMermaid図を生成
            mermaid_contents = []
            for chunk_data in chunks:
                mermaid_writer = WriteMermaid(chunk_data)
                mermaid_content = mermaid_writer.generate()
                mermaid_contents.append(mermaid_content)
            
            # ファイルに出力
            file_writer = WriteFile(self.output_file)
            
            # 出力が1つだけならそのまま書き込み、複数なら分割書き込み
            if len(mermaid_contents) == 1:
                file_writer.write(mermaid_contents[0])
                print(f"解析が完了しました。出力ファイル: {self.output_file}")
            else:
                file_writer.write_split(mermaid_contents)
                print(f"解析が完了しました。{len(mermaid_contents)}個のファイルに分割出力しました。")
            
            return True
            
        except Exception as e:
            print(f"解析中にエラーが発生しました: {str(e)}")
            import traceback
            traceback.print_exc()
            return False
