#!/usr/bin/env python3
"""
PCAPファイル読み込みクラスを提供するモジュール
"""
import pyshark


class ReadPcap:
    """
    PCAPファイルを読み込むクラス
    """
    
    def __init__(self, pcap_file):
        """
        ReadPcapクラスの初期化
        
        Args:
            pcap_file (str): PCAPファイルのパス
        """
        self.pcap_file = pcap_file
    
    def read(self):
        """
        PCAPファイルを読み込み、パケットのリストを返す
        
        PySharkを使用してPCAPファイルを読み込む
        
        Returns:
            list: pyshark.packet.Packetオブジェクトのリスト
        """
        try:
            print(f"PCAPファイル読み込み中: {self.pcap_file}")
            # PySharkを使用してPCAPファイルを読み込む
            capture = pyshark.FileCapture(self.pcap_file)
            
            # パケットのリストに変換
            packets = list(capture)
            print(f"読み込み完了: {len(packets)}パケット")
            
            return packets
            
        except Exception as e:
            print(f"PCAPファイル読み込み中にエラーが発生しました: {str(e)}")
            import traceback
            traceback.print_exc()
            raise
