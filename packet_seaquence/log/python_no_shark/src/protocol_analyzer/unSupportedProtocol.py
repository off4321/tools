from src.protocol_analyzer.base import ProtocolAnalyzer

class AnalyzeUnsupportedProtocol(ProtocolAnalyzer):
    """
    サポートされていないプロトコルの解析クラス
    """
    
    def __init__(self, packet, highest_layer="Unknown"):
        """
        AnalyzeUnsupportedProtocolクラスの初期化
        
        Args:
            packet: scapy.packet.Packetオブジェクト
            highest_layer (str): 最上位レイヤー名
        """
        super().__init__(packet)
        self.highest_layer = highest_layer
    
    def analyze(self):
        """
        未サポートプロトコル情報を解析する
        
        Returns:
            dict: プロトコル名を含む辞書
        """
        return {
            "protocol_name": self.highest_layer,
            "supported": False
        }
    
    def get_display_info(self):
        """
        表示用の情報を取得する
        
        Returns:
            str: 表示用のプロトコル情報
        """
        return f"{self.highest_layer}[未サポート]"
    
    def get_summary(self):
        """
        プロトコル情報のサマリーを取得する
        
        Returns:
            str: プロトコル情報のサマリー
        """
        return f"{self.highest_layer}[未サポート]"
