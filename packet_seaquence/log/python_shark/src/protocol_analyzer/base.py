#!/usr/bin/env python3
"""
プロトコル解析クラスの基底クラスを提供するモジュール
"""


class ProtocolAnalyzer:
    """
    プロトコル解析クラスの基底クラス
    
    各プロトコル解析クラスはこのクラスを継承して実装します
    """
    
    def __init__(self, packet):
        """
        ProtocolAnalyzerクラスの初期化
        
        Args:
            packet: pyshark.packet.Packetオブジェクト
        """
        self.packet = packet
    
    def analyze(self):
        """
        パケットからプロトコル情報を解析する
        
        Returns:
            dict: プロトコル情報を含む辞書
        """
        raise NotImplementedError("このメソッドはサブクラスでオーバーライドする必要があります")
    
    def get_display_info(self):
        """
        表示用の情報を取得する
        
        Returns:
            str: 表示用のプロトコル情報
        """
        raise NotImplementedError("このメソッドはサブクラスでオーバーライドする必要があります")
    
    def get_summary(self):
        """
        プロトコル情報のサマリーを取得する
        
        Returns:
            str: プロトコル情報のサマリー
        """
        raise NotImplementedError("このメソッドはサブクラスでオーバーライドする必要があります")
