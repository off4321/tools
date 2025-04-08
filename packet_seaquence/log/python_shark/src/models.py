#!/usr/bin/env python3
"""
パケットシーケンス解析に必要なデータモデルを定義するモジュール
"""

class PacketSequenceData:
    """
    パケットシーケンス情報を含むデータ構造
    解析されたパケット情報を保存するためのクラス
    """
    def __init__(self):
        """
        PacketSequenceDataクラスの初期化
        """
        self.packets = []  # パケット情報のリスト
        self.nodes = set()  # シーケンス図に表示するノード（IPアドレスやMACアドレスなど）
    
    def add_packet(self, packet_info):
        """
        解析されたパケット情報を追加する
        
        Args:
            packet_info (dict): パケット情報を含む辞書
                                {
                                    'time': パケットキャプチャ時刻,
                                    'src': 送信元アドレス,
                                    'dst': 宛先アドレス,
                                    'protocol': プロトコル名,
                                    'info': 詳細情報の辞書
                                }
        """
        self.packets.append(packet_info)
        self.nodes.add(packet_info['src'])
        self.nodes.add(packet_info['dst'])
    
    def get_packets(self):
        """
        解析されたパケット情報のリストを取得する
        
        Returns:
            list: パケット情報のリスト
        """
        return self.packets
    
    def get_nodes(self):
        """
        シーケンス図に表示するノードのセットを取得する
        
        Returns:
            set: ノードのセット
        """
        return self.nodes
    
    def copy_first_n(self, n):
        """
        最初のn個のパケットだけをコピーした新しいPacketSequenceDataオブジェクトを返す
        
        Args:
            n (int): コピーするパケット数
            
        Returns:
            PacketSequenceData: 新しいパケットシーケンスデータオブジェクト
        """
        new_data = PacketSequenceData()
        
        # 最初のn個のパケットをコピー
        packets = self.get_packets()
        for i, packet in enumerate(packets):
            if i >= n:
                break
            new_data.add_packet(packet)
        
        return new_data
