#!/usr/bin/env python3
"""
PCAPファイル読み込みクラスを提供するモジュール
"""
from scapy.all import rdpcap


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
        
        Scapyを使用してPCAPファイルを読み込む
        
        Returns:
            list: scapy.packet.Packetオブジェクトのリスト
        """
        try:
            print(f"PCAPファイル読み込み中: {self.pcap_file}")
            # Scapyを使用してPCAPファイルを読み込む
            packets = rdpcap(self.pcap_file)
            
            print(f"読み込み完了: {len(packets)}パケット")
            
            return packets
            
        except Exception as e:
            print(f"PCAPファイル読み込み中にエラーが発生しました: {str(e)}")
            import traceback
            traceback.print_exc()
            raise
