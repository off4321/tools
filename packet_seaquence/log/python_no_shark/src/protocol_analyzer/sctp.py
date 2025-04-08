from scapy.all import SCTP
from src.protocol_analyzer.base import ProtocolAnalyzer

class AnalyzeSctp(ProtocolAnalyzer):
    """
    SCTP内の詳細情報を解析するクラス
    """

    def analyze(self):
        if SCTP not in self.packet:
            return None
        try:
            sctp_info = {}
            sctp_layer = self.packet[SCTP]
            
            # 送信元ポート
            sctp_info['src_port'] = str(sctp_layer.sport)
            
            # 宛先ポート
            sctp_info['dst_port'] = str(sctp_layer.dport)
            
            # 検証タグ
            if hasattr(sctp_layer, 'tag'):
                sctp_info['verification_tag'] = str(sctp_layer.tag)
                
            # チェックサム
            if hasattr(sctp_layer, 'chksum'):
                sctp_info['checksum'] = hex(sctp_layer.chksum) if sctp_layer.chksum is not None else "None"
            
            return sctp_info
            
        except Exception as e:
            print(f"SCTP解析中にエラーが発生しました: {str(e)}")
            import traceback
            traceback.print_exc()
            return None

    def get_display_info(self):
        info = self.analyze()
        if not info:
            return "SCTP情報なし"
        return f"SCTP {info.get('src_port','?')} -> {info.get('dst_port','?')} (VTAG={info.get('verification_tag','?')})"

    def get_summary(self):
        info = self.analyze()
        if not info:
            return "SCTP"
        return f"SCTP {info.get('src_port','?')} -> {info.get('dst_port','?')}"